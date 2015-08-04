package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"webo/models/rpc"
	"webo/models/svc"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	//    role := this.GetSession("role")
	//    if role == nil {
	this.Redirect("/static/frame/login.html", 302)
	//    } else {
	//        this.Redirect("/static/main.html", 302)
	//    }
}
func (this *LoginController) Post() {
	username := this.GetString("login_username")
	password := this.GetString("login_password")

	loginRet := rpc.JsonResult{}
	if username == "" || password == "" {
		loginRet.Result = "请输入用户名和密码！"
	}
	params := svc.Params{
		"username":username,
		"password":password,
	}
	code, user := svc.Get("user", params)

	if code != "success" {
		fmt.Println("err", code)
		loginRet.Result = "用户名或密码错误！"
	} else {

		this.SetSession("username", username)
		if role, ok:= user["role"];ok{
			this.SetSession("role", role)
			loginRet.Ret = "success"
		}else {
			loginRet.Ret = "faild"
			loginRet.Result = "获取权限失败"
		}
	}
	this.Data["json"] = &loginRet
	this.ServeJson()
}

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Get() {
	this.DelSession("role")
	this.DelSession("username")
	this.Redirect("/", 302)
}
