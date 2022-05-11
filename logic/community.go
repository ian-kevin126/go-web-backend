package logic

import (
	"gin_demo/dao/mysql"
	"gin_demo/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库，获取社区列表并返回
	return mysql.GetCommunityList()
}

/**
 * @Author huchao
 * @Description //TODO 根据ID查询分类社区详情
 * @Date 17:08 2022/2/12
 **/
func GetCommunityDetailByID(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityByID(id)
}
