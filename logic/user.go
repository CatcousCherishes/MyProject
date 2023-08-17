package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

//	logic/service：逻辑（服务）层，负责处理业务逻辑,  存放业务逻辑处理代码
//
// 1.SignUp()这个函数需要做 注册业务逻辑处理
func SignUp(p *models.ParamSignUp) (err error) {
	//1. 判断用户不存在 （这里是注册用户，如果存在就不继续了）
	if err = mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询出错
		return err
	}

	//继续下面 用户不存在，则开始注册,生成user_ID，并保存到数据库
	//2.生成UID (雪花算法生成ID)
	userID := snowflake.GenID()
	//构造一个User实例  user
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.保存到数据库 DAO层
	return mysql.IntertUser(user)
}

// 2.Login 登录业务流程处理
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//直接执行用户登入

	//传递的是指针，就能拿到 user.userID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// username password 校验正确 登录成功后，生成JWT的token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return user, err
}
