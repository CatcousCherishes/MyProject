package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
)

// 社区相关的---

// CommunityHandler 社区列表
// @Summary 社区列表
// @Description 社区列表
// @Tags 社区业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.Community false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	//遇事不决，先写注释
	//1.查询到所有社区（community_id,community_name） 以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错 返回给前端
		return
	}
	ResponseSuccess(c, data)
}

// 社区分类详情
// CommunityDetailHandler 社区详情
// @Summary 社区详情
// @Description 社区详情
// @Tags 社区业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query  models.CommunityDetail  false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /community/:id [get]
func CommunityDetailHandler(c *gin.Context) {
	//1.从前端请求中  获取社区id(拿到参数详细步骤)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.根据id 获取社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错 返回给前端
		return
	}
	ResponseSuccess(c, data)
}
