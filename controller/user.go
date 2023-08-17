package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"
)

// Controller:服务的入口，负责处理路由、参数校验、请求转发
// 规矩：添加代码 从请求开始，按处理流程走 前端VUE - NGINX -HTTP/Thrift/gRPC  - Controller - Logic - DAO

// 功能1、用户注册 SignUpHander 处理注册请求的函数
// SignUpHandler 注册业务
// @Summary 注册业务
// @Description 注册业务
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /signup [POST]
func SignUpHander(c *gin.Context) {
	//1. 获取请求参数和参数校验
	//var p models.ParamSignUp
	//ShouldBindJSON只能判断参数是否json
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil { //从c 中通过ShouldBindJSON获取到 p
		//请求参数有误,直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err)) //zap记录错误信息，直接zap.Error()
		errs, ok := err.(validator.ValidationErrors)               //判断err是不是validator.ValidationErrors 类型
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans))) //翻译并除掉错误提示中的结构体标识
		return
	}
	// 打印参数p
	fmt.Println(p) // 已经拿到 请求参数 p变量

	//2.业务处理 放到logic层,具体在调用
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应给用户（成功）
	ResponseSuccess(c, nil)
}

// 功能2、LoginHander 用户登录函数 处理登入请求的函数
// LoginHandler 登录业务
// @Summary 登录业务
// @Description 登录业务
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamLogin false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /login [POST]
func LoginHandler(c *gin.Context) {
	//1. 获取参数和校验参数
	p := new(models.ParamLogin)
	//从c 中通过ShouldBindJSON获取到 p
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误,直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err)) //zap记录错误信息，直接zap.Error()
		//判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err)) //出错记录
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			//中断失败
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3.返回响应，带上token,应当是返回user
	//id json与int类型之间会存在失真的情况 P57
	// 序列化：后端准备数据的时候 int64,发送给前端需要转换成string类型；前端传值是字符串类型的，也必须解析成int64，后端才能接收到int64，

	// 当生成的id> json 的最大范围，id值大于1<<53-1  int64类型的最大值是1<<63-1；json中无法表示64位的数据，
	//如何解决：因后端需要将其转换为字符串类型。

	//记住理解查询 go语言操作JSON的技巧？
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), // 将UserID int类型，转换成string ，返回到前端
		"user_name": user.Username,
		"token":     user.Token,
	})
}
