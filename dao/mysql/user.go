package mysql

import (
	"blog_app/model"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

var Secret = "cisco46589"

func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username=?`
	var count int64
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

func InsertUser(user *model.User) (err error) {
	//1、密码加密
	user.Password = encryptPassword(user.Password)

	//2、写入库
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserId, user.Username, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(Secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *model.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id,username,password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}

	if user.Password != encryptPassword(oPassword) {
		return ErrorInvalidPassword
	}
	return

}

func GetUserById(uid int64) (user *model.User, err error) {
	user = new(model.User)
	sqlStr := "select user_id,username from user where user_id=?"
	err = db.Get(user, sqlStr, uid)
	return
}
