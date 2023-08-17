package redis

//设计 redis key

//redis key 需要注意 使用命名空间的方式，方便查询和拆分

const (
	Prefix                 = "bulebell:"   //项目key前缀
	KeyPostTimeZSet        = "post:time"   // zset;按帖子及发帖时间  bulebell:post:time
	KeyPostScoreZSet       = "post:score"  //zset；按帖子与帖子的得分 bulebell:Post:score
	KeyPostVotedZSetPrefix = "post:voted:" // zest；记录用户及投票类型;参数是post_id（不完整的key）  bulebell:Post:score  bulebell:post:voted:post_id

	KeyCommunityPostSetPrefix = "bluebell:community:" // set保存每个分区下帖子的id
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
