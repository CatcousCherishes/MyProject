package models

// 用来定义请求参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required" `
	Password   string `json:"password" binding:"required" `
	RePassword string ` json:"re_Password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	//UserID 可以从当前登录的用户中获取到
	PostID    string `json:"post_id" binding:"required" `             //帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票（1）反对票（-1） 取消投票（0）,限定只能是1 0 -1
}

// ParamPostList 查询获取帖子列表query string 参数
type ParamPostList struct {
	CommunityID uint64 `json:"community_id",form:"community_id"` //可以为空
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
}

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamCommunityPostList
type ParamCommunityPostList struct {
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
	CommunityID int64  `json:"community_id",form:"community_id"`
}
