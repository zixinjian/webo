package ctrl

import (
	"github.com/astaxie/beego"
	"wb/cc"
	"wb/cs"
	"wb/lg"
	"wb/st"
	"wb/usr"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) LoginAs(u usr.Usr) {
	this.SetSession(cc.SessionUser, u)
	lg.SInfo("[S] login as:", u.GetUserName(), "role:", u.GetRole())
	loginRet := cs.JsonResult{}
	loginRet.Ret = st.Success
	loginRet.Result = "登录成功！"
	this.Data["json"] = &loginRet
	this.ServeJSON()
}

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Get() {
	this.DestroySession()
	this.Redirect("/", 302)
}
