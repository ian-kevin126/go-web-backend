package logic

import (
	"gin_demo/dao/mysql"
	"gin_demo/dao/redis"
	"gin_demo/models"
	"gin_demo/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	// 1、 生成post_id(生成帖子ID)
	post.PostID = snowflake.GenID()

	// 2、创建帖子 保存到数据库
	err = mysql.CreatePost(post)
	if err != nil {
		return err
	}

	err = redis.CreatePost(post.PostID, post.CommunityID)
	return
}

func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想要的数据格式
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}

	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorId) failed", zap.Int64("post.AuthorId", post.AuthorId), zap.Error(err))
		return
	}

	// 根据社区ID查询社区信息
	community, err := mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID(post.CommunityID) failed", zap.Int64("post.CommunityID", post.CommunityID), zap.Error(err))
		return
	}

	data = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		Post:            post,
		CommunityDetail: community,
	}

	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	list, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, len(list))

	for _, post := range list {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId) failed", zap.Int64("post.AuthorId", post.AuthorId), zap.Error(err))
			continue
		}

		// 根据社区ID查询社区信息
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID(post.CommunityID) failed", zap.Int64("post.CommunityID", post.CommunityID), zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)
	}

	return
}

func GetPostListByCreateTimeOrScore(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2，去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}

	zap.L().Debug("GetPostListByCreateTimeOrScore", zap.Any("ids", ids))
	// 3，根据id去数据库查询帖子的详细信息
	// 返回的数据还要按照我给定的id的顺序返回（通过在sql中加入：order by FIND_IN_SET(post_id, ?)）
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId) failed", zap.Int64("post.AuthorId", post.AuthorId), zap.Error(err))
			continue
		}

		// 根据社区ID查询社区信息
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID(post.CommunityID) failed", zap.Int64("post.CommunityID", post.CommunityID), zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: community,
			VoteNum:         voteData[idx],
		}

		data = append(data, postDetail)
	}

	return
}

/**
 * @Author huchao
 * @Description //TODO  根据社区去查询帖子列表
 * @Date 22:53 2022/2/16
 **/
func GetCommunityPostList(post *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2、去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(post)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostList(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3、根据id去数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回  order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Int64("postID", post.AuthorId),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
			AuthorName:      user.UserName,
		}
		data = append(data, postDetail)
	}
	return
}

/**
 * @Author huchao
 * @Description //TODO 将两个查询帖子列表逻辑合二为一的函数
 * @Date 12:08 2022/2/17
 **/
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同,执行不同的业务逻辑
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostListByCreateTimeOrScore(p)
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
