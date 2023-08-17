package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"web_app/logic"
	"web_app/models"
)

//投票功能函数

//type VoteData struct {
//	//UserID 可以从当前登录的用户中获取到
//	PostID    int64 `json:"post_id,string"`   //帖子id
//	Direction int   `json:"direction,string"` // 赞成票（1）反对票（-1）
//}

// PostVoteController 投票
// @Summary 投票
// @Description 投票
// @Tags 投票业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /vote [POST]
func PostVoteController(c *gin.Context) {
	//参数校验，给哪个帖子文章投什么票
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// if 参数时符合validator情况的
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData) //翻译并除掉错误提示中的结构体标识
		return
	}
	//获取当前用户的ID  getCurrentUser
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//logic 具体投票的逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回成功响应
	ResponseSuccess(c, nil)
}
