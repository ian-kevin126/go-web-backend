package redis

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

/*
投票的几种情况：
direction=1时，有两种情况：
	1，之前没有投过票，现在投赞成票 ---> 更新分数和投票记录 差值：1 +432
	2，之前投反对票，现在改投赞成票 ---> 更新分数和投票记录 差值：2 +432 * 2

direction=0，有两种情况：
	1，之前投反对票，现在要取消投票 ---> 更新分数和投票记录 差值：1 +432
	2，之前投赞成票，现在要取消投票 ---> 更新分数和投票记录 差值：1 -432

direction=-1，有两种情况：
	1，之前没有投过票，现在投反对票 ---> 更新分数和投票记录 差值：1 -432
	2，之前投赞成票，现在改投反对票 ---> 更新分数和投票记录 差值：2 -432 * 2

规律：如果现在的值大于之前的值，那就是+432

投票限制：每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1，到期之后将redis中保存的赞成票及反对票存储到mysql表中
	2，到期之后删除那个 KeyPostVotedZSetPrefix
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值的分数
	PostPerAge       = 20
)

func VoteForPost(userID, postID string, value float64) (err error) {
	// 1，判断投票的限制
	// 去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	fmt.Println("投票时间：", float64(time.Now().Unix())-postTime)
	// 大于一周投票就过期
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}

	// 2，更新帖子的分数
	// 2和3 需要放到一个pipeline事务中操作
	// 判断是否已经投过票 查当前用户给当前帖子的投票记录
	// 先查询当前用户给当前帖子的投票记录
	key := KeyPostVotedZSetPrefix + postID
	oldValue := client.ZScore(getRedisKey(key), userID).Val()

	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if oldValue == value {
		return ErrorVoteRepeat
	}

	var op float64
	if value > oldValue {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(oldValue - value)                                                              // 计算两次投票的差值
	pipeline := client.TxPipeline()                                                                 // 事务操作
	_, err = pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), scorePerVote*diff*op, postID).Result() // 更新分数
	if ErrorVoteTimeExpire != nil {
		return err
	}

	// 3、记录用户为该帖子投票的数据
	if value == 0 {
		_, err = client.ZRem(getRedisKey(key), postID).Result()
	} else {
		pipeline.ZAdd(getRedisKey(key), redis.Z{ // 记录已投票
			Score:  value, // 赞成票还是反对票
			Member: userID,
		})
	}

	_, err = pipeline.Exec()

	return err
}

func CreatePost(postID, communityID int64) error {
	// 事务操作
	pipeline := client.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunityPostSetPrefix + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)

	_, err := pipeline.Exec()

	return err
}
