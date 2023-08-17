package redis

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"web_app/models"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3.ZREVRANGE 按照分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

func GetPostIDInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis 中获取id的3步
	//究竟取哪一种zset的id序列，需要根据time和score
	// 1. 根据用户请求中携带的order参数来确定查询的redis key
	key := getRedisKey(KeyPostTimeZSet) //默认是按照time来排序  取redis key (KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.确定查询的索引的起始点
	start := (p.Page - 1) * p.Size //第一页从0开始
	end := start + p.Size - 1
	//3.redis ZRevRange 按分数从大到小的顺序查询指定数量的元素 ids []string
	return client.ZRevRange(key, start, end).Result()
}

// GetPostVoteData() 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPrefix + id)
	//	//查找key元中分数是1的元素的数量->也就是统计每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val() //统计一下 1 的数量
	//	data = append(data, v)
	//}
	//return data, err

	// 使用 pipeline一次发送多条命令,减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// TODO 按社区查询ids(查询出的ids已经根据order从大到小排序)
func GetCommunityPostIDsInOrder(p *models.ParamCommunityPostList) ([]string, error) {
	// 1.根据用户请求中携带的order参数确定要查询的redis key
	orderkey := KeyPostTimeZSet       // 默认是时间
	if p.Order == models.OrderScore { // 按照分数请求
		orderkey = KeyPostScoreZSet
	}

	// 使用zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据

	// 社区的key
	cKey := KeyCommunityPostSetPrefix + strconv.Itoa(int(p.CommunityID))

	// 利用缓存key减少zinterstore执行的次数 缓存key
	key := orderkey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX", // 将两个zset函数聚合的时候 求最大值
		}, cKey, orderkey) // zinterstore 计算
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在的就直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
