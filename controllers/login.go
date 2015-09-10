package controllers

import (
	//	"fmt"
	"github.com/astaxie/beego"
	"webo/models/rpc"
	"webo/models/svc"
	"webo/models/t"
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
	params := t.Params{
		"username": username,
		"password": password,
	}
	code, user := svc.Get("user", params)

	if code != "success" {
		beego.Error("User:%s login error code: %s", username, code)
		loginRet.Result = "用户名或密码错误！"
		this.Data["json"] = &loginRet
		this.ServeJson()
		return
	}
	loginRet = this.setSessionFromUser(user)
	if loginRet.Ret == "success" {
		this.SetSession(SessionUserName, username)
	}
	this.Data["json"] = &loginRet
	this.ServeJson()
}
func (this *LoginController) setSessionFromUser(user map[string]interface{}) rpc.JsonResult {
	loginRet := rpc.JsonResult{}
	role, ok := user["role"]
	if !ok {
		loginRet.Ret = "faild"
		loginRet.Result = "获取权限失败"
		return loginRet
	}
	sn, ok := user["sn"]
	if !ok {
		loginRet.Ret = "faild"
		loginRet.Result = "获取权限失败"
		return loginRet
	}
	department, ok := user["department"]
	if !ok {
		loginRet.Ret = "faild"
		loginRet.Result = "获取权限失败"
		return loginRet
	}
	this.SetSession(SessionUserRole, role)
	this.SetSession(SessionUserSn, sn)
	this.SetSession(SessionUserDepartment, department)
	loginRet.Ret = "success"
	loginRet.Result = "登录成功"
	return loginRet
}

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Get() {
	this.DelSession(SessionUserName)
	this.DelSession(SessionUserRole)
	this.DelSession(SessionUserDepartment)
	this.DelSession(SessionUserSn)
	this.Redirect("/", 302)
}
