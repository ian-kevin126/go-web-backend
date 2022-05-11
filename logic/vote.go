package logic

import (
	"gin_demo/dao/redis"
	"gin_demo/models"
	"strconv"

	"go.uber.org/zap"
)

// 投票功能，推荐阅读（阮一峰-基于用户投票相关的算法）：https://www.ruanyifeng.com/blog/algorithm/
// 本项目使用简化版的投票分数，投一票就加432分，86400/200 ——》需200张赞成票可以给你的帖子续一天（例子：来自redis实战）

/*
投票的几种情况：
direction=1时，有两种情况：
	1，之前没有投过票，现在投赞成票 ---> 更新分数和投票记录
	2，之前投反对票，现在改投赞成票 ---> 更新分数和投票记录

direction=0，有两种情况：
	1，之前投赞成票，现在要取消投票 ---> 更新分数和投票记录
	2，之前投反对票，现在要取消投票 ---> 更新分数和投票记录

direction=-1，有两种情况：
	1，之前没有投过票，现在投反对票 ---> 更新分数和投票记录
	2，之前投赞成票，现在改投反对票 ---> 更新分数和投票记录

投票限制：每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1，到期之后将redis中保存的赞成票及反对票存储到mysql表中
	2，到期之后删除那个 KeyPostVotedZSetPrefix
*/

func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}