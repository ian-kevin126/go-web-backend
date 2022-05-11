package controllers

import (
	"gin_demo/logic"
	"gin_demo/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

/**
 * @Author ian-kevin
 * @Description //TODO 创建帖子
 * @Date 17:40 2022/2/12
 **/
// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.Post false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /post [POST]
func CreatePostHandler(c *gin.Context) {
	// 1，获取参数及参数的校验
	// c.shouldBindJSON gin内部会去调用validator，type中的binding tag就是这个实现
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 从context获取当前发请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	p.AuthorId = userID

	// 2，创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3，返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	// 1，获取参数（帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2，根据id取出帖子的数据
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3，返回响应
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	page, size := getPageInfo(c)
	// 1，获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 2，返回响应
	ResponseSuccess(c, data)
}

/*
GetPostListByCreateTimeOrScoreHandler
根据前端传来的参数动态获取帖子列表
根据创建时间或者分数获取帖子列表

步骤：
1，获取参数
2，去redis查询id列表
3，根据id去数据库查询帖子详细信息
*/
func GetPostListByCreateTimeOrScoreHandler(c *gin.Context) {
	// 1，获取参数
	// GET请求-/api/v1/post/lists?page=1&size=10&order=time  queryString参数
	// 获取分页参数
	// c.ShouldBind()：根据请求的数据类型选择相应的方法取获取数据
	// c.ShouldBindQuery()：如果请求中携带的是json格式的数据，才能用这个方法获取到数据

	// 初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListByCreateTimeOrScoreHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2，去redis查询id数据
	data, err := logic.GetPostListByCreateTimeOrScore(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3，返回响应
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler 根据社区去查询帖子列表
func GetCommunityPostListHandler(c *gin.Context) {
	// GET请求参数(query string)： /api/v1/posts2?page=1&size=10&order=time
	// 获取分页参数
	p := &models.ParamPostList{
		Page:        1,
		Size:        10,
		Order:       models.OrderScore,
		CommunityID: 0,
	}
	//c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 获取数据
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
