package logic

import (
	"blog_app/dao/mysql"
	"blog_app/model"
	"blog_app/pkg/jwt"
	"blog_app/pkg/snowflake"
)

func SignUp(p *model.ParamSignUp) (err error) {
	//1、判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	//2、生成uid
	uid := snowflake.GenID()

	//3、添加到数据库
	user := &model.User{
		UserId:   uid,
		Username: p.Username,
		Password: p.Password,
	}

	return mysql.InsertUser(user)

}

func Login(p *model.ParamLogin) (user *model.User, err error) {

	//1、调用登录
	user = &model.User{
		Username: p.Username,
		Password: p.Password,
	}

	if err = mysql.Login(user); err != nil {
		return nil, err
	}

	//2、生成token

	token, err := jwt.GenToken(user.UserId, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return

}
