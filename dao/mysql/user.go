package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"web_app/models"
)

// 把每一步数据库操作封装成函数
// 待logic 层根据业务需求调用
const screct = "liwenzhou.com"

// CheckUserExist 检验指定用户是否已经存在
func CheckUserExist(username string) (err error) {
	//数据库执行查询语句
	sqlstr := `select count(user_id) from user where username = ? `
	var count int
	if err := db.Get(&count, sqlstr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func IntertUser(user *models.User) (err error) {
	//需要对密码进行加密后，再存放到数据库
	user.Password = encryptPassword(user.Password)
	//执行SQL语句入库
	sqlstr := `insert into user(user_id,username,password) values (?,?,?)`
	_, err = db.Exec(sqlstr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(screct)) //加严 字符串
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// login 登录
func Login(user *models.User) (err error) {
	oPassword := user.Password // 记录用户原始登录的密码
	//查询用户
	sqlstr := `select user_id,username,password from user where username=?`
	if err := db.Get(user, sqlstr, user.Username); err != nil {
		return err
	}
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		//查询数据库失败
		return err
	}
	//判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// 根据 (用户)作者id   获取作者信息
func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}
