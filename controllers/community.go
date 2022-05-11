package controllers

import (
	"gin_demo/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CommunityHandler(c *gin.Context) {
	// 1，查询到所有的社区（community_id,community_name），以列表的形式返回
	data, err := logic.GetCommunityList()

	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed.", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外界
		return
	}

	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	// 1，获取社区id
	communityID := c.Param("id")
	// 因为communityID是string类型的，先将id转化成十进制int类型
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2，查询到社区详情
	data, err := logic.GetCommunityDetailByID(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}
