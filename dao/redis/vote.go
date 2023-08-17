package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

// 本项目使用简化版的投票分数
//《redis实战》书里面的经典例子
// 投一票就加432分   86400s(一天86400秒)/200(一天获得200票)  --> 200张赞成票可以给你的帖子续一天

/*投票的几种情况
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票    --> 更新分数和投票记录   差值的绝对值为：1  +432
	2. 之前投反对票，现在改投赞成票    --> 更新分数和投票记录   差值的绝对值为：2  +432*2
direction=0时，有两种情况：
    1. 之前投过反对票，现在要取消投票  --> 更新分数和投票记录   差值的绝对值为：1  +432
（上面是：现在的值 vaule > ov 原来的值 ）

	2. 之前投过赞成票，现在要取消投票  --> 更新分数和投票记录   差值的绝对值为：1  -432

direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票    --> 更新分数和投票记录   差值的绝对值为：1  -432
	2. 之前投赞成票，现在改投反对票    --> 更新分数和投票记录   差值的绝对值为：2  -432*2

投票的限制：
每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2. 到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepested   = errors.New("不能重复投票")
)

func CreatePost(postID, CommunityID int64) error {
	//这两个事务必须要同时成功，使用pipeline
	pipeline := client.TxPipeline()

	//保存帖子创建(发帖)的时间到 rediskey -  KeyPostTimeZSet
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子分数：是在当前的时间啊基础上 +或者-432 (也就是帖子时间是基准)
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 更新： 把帖子id 加到社区set里面
	communityKey := KeyCommunityPostSetPrefix + strconv.Itoa(int(CommunityID))
	pipeline.SAdd(communityKey, postID) // 添加到对应版块  把帖子添加到社区的set

	_, err := pipeline.Exec()
	return err

}

//go操作redis数据库是有哪些技巧

// VoteForPost 帖子投票
func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票的限制(可以从帖子的时间着手)
	//去redis取出帖子的发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2.更新帖子的分数

	//2和3需要放到一个pipeline
	pipeline := client.TxPipeline()
	//先查询之前的投票记录,就是先查 当前用户给当前帖子的投票记录,拿到一个分数
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	//更新：如果这一次的投票值与之前保存的值一致时，则是重复投票
	if value == ov {
		return ErrVoteRepested
	}
	//这个是在找方向，>就是+，< 就是-
	var dir float64
	if value > ov {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的差值
	//Redis Zincrby 命令  是把Golang postID的分数加dir*diff*432
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), dir*diff*scorePerVote, postID)
	//更新分数是否成功

	//3.记录用户为帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value, //赞成还是反对
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
