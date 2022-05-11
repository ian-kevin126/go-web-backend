package redis

import "errors"

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoted          = errors.New("已经投过票了")
	ErrorVoteRepeat     = errors.New("请勿重复投票")
)
