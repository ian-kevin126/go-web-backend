package mysql

import (
	"database/sql"
	"gin_demo/models"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"

	if err := db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("There is no community in db.")
			err = nil
		}
	}

	return
}

/**
 * @Author ian-kevin
 * @Description // TODO 根据ID查询分类社区详情
 * @Date 17:08 2022/2/12
 **/
func GetCommunityByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time
	from community
	where community_id = ?`
	err = db.Get(community, sqlStr, id)

	if err == sql.ErrNoRows { // 查询为空
		err = ErrorInvalidID // 无效的ID
		return
	}

	if err != nil {
		zap.L().Error("query community failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
	}

	return community, err
}
