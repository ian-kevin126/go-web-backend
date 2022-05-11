package models

import "time"

/**
 * @Author huchao
 * @Description //TODO 帖子Post结构体
 * @Date 17:44 2022/2/12
 **/
// 内存对齐概念 字段类型相同的对齐 缩小变量所占内存大小
type Post struct {
	PostID      int64     `json:"post_id,string" db:"post_id"`
	AuthorId    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"-" db:"create_time"`
}

/**
 * @Author ian-kevin
 * @Description //TODO 帖子返回的详情结构体
 * @Date 21:59 2022/2/12
 **/
type ApiPostDetail struct {
	*Post                               // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息
	AuthorName       string             `json:"author_name"`
	VoteNum          int64              `json:"vote_num"`
	//CommunityName string `json:"community_name"`
}

/*
 * ParamPostList
 * @Author ian-kevin
 * @Description //TODO 获取帖子列表query string参数
 * @Date 21:51 2022/2/15
 **/
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page"`                   // 页码
	Size        int64  `json:"size" form:"size"`                   // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}
