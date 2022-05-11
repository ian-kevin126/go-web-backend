package redis

import (
	"gin_demo/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

/**
 * @Author huchao
 * @Description //TODO 按照分数从大到小的顺序查询指定数量的元素
 * @Date 0:12 2022/2/17
 **/
func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3.ZREVRANGE 按照分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 1，根据用户请求中携带的order参数确定要插叙的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	// 查询redis指令：redis.cn
	// 2，确定查询的索引起点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1

	// 3，ZRevRange查询，按分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

/**
 * @Author huchao
 * @Description //TODO 根据ids查询每篇帖子的投赞成票的数据
 * @Date 21:28 2022/2/16
 **/
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids{
	//	key := KeyPostVotedZSetPrefix + id
	//	// 查找key中分数是1的元素数量 -> 统计每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}

	// 使用 pipeline 一次发送多条命令，减少RTT
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
		// 转换格式
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}

	return
}

/**
 * @Author huchao
 * @Description //TODO 按社区查询ids(查询出的ids已经根据order从大到小排序)
 * @Date 23:06 2022/2/16
 * @Param orderKey:按照分数或时间排序
	将社区key与orderkey(社区或时间)做zinterstore
 **/
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 1.根据用户请求中携带的order参数确定要查询的redis key
	orderKey := KeyPostTimeZSet       // 默认是时间
	if p.Order == models.OrderScore { // 按照分数请求
		orderKey = KeyPostScoreZSet
	}

	// 使用zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据

	// 社区的key
	cKey := getRedisKey(KeyCommunityPostSetPrefix + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少zinterstore执行的次数 缓存key
	key := orderKey + strconv.Itoa(int(p.CommunityID))

	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX", // 将两个zset函数聚合的时候 求最大值
		}, cKey, orderKey) // zinterstore 计算

		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}

	// 存在的就直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
