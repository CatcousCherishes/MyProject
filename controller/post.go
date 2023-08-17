package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
	"web_app/models"
)

// CreatePostHandler 创建帖子的处理函数
// 1.创建数据库中的post表结构
// 2.定义好post对应的模型结构体
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
	//1.获取参数及参数校验
	//c.ShouldBindJSON()  // validator --> binding tag
	p := new(models.Post)
	//从c 中 通过ShouldBindJSON获取到 p
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从c 取到当前发送请求的用户的ID
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p)  failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 功能 查询单个获取帖子详情
// @Summary 获取单个帖子
// @Description 获取帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts [GET]
func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数 （从URL中获取帖子的id ） /api/v1/post/:id
	pidStr := c.Param("id")
	//转换成Int类型,利用strconv.ParseInt 将字符串解析成Int
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据取出的id ,查询数据库中帖子的数据
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应 带上查询结果 data//优化帖子详情 先修改返回数据data的字段
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表[]
// @Summary 分页获取帖子列表
// @Description 分页获取帖子列表
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts [GET]
func GetPostListHandler(c *gin.Context) {
	//获取分页参数 处理 getPageInfo函数获取到page\size
	page, size := getPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// 升级版 GetPostListHandler2  根据时间和分数来获取post列表
// 根据前端 传入的参数（按创建时间、或者按照帖子的分数 排序）
// 1、获取参数
// 2、去redis 查询id列表
// 3、根据id去数据库mysql查询帖子详情信息
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	//Get 前端请求参数(query string) : /api/v1/post2?page=1&size=10&order=time
	// 初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, //通过常数 const 避免出现magic string,默认按照时间排序
	}
	//c.ShouldBind() 根据请求的数据类型 选择相应的方法去获取到数据
	//c.ShouldBindJSON()  如果请求中携带的数据是json格式数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with Invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// 根据社区去查询帖子列表

func GetCommunityPostListHandler(c *gin.Context) {
	// 初始化结构体时指定初始参数
	p := &models.ParamCommunityPostList{
		CommunityID: 1,
		Page:        1,
		Size:        10,
		Order:       models.OrderTime,
	}
	//c.ShouldBind()  根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取数据
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
