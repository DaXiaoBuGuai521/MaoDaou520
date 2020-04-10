package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
)

type UserController struct {
	beego.Controller
}

//展示注册页面
func (this *UserController) ShowRegister() {

	this.TplName = "register.html"
}

//处理注册页面
func (this *UserController) HandleReg() {
	// 接受数据

	userName := this.GetString("userName")
	pwd := this.GetString("password")
	//校验数据
	if userName == "" || pwd == "" {
		beego.Error("用户名或密码为空")
		//返回前段数据
		this.Data["errmsg"] = "用户名或密码不能为空"
		return
	}

	//处理数据
	//插入数据
	o := orm.NewOrm()

	//插入对象
	var user models.User
	//给插入对象赋值
	user.UserName = userName
	user.Pwd = pwd

	//插入
	_, err := o.Insert(&user)
	if err != nil {
		this.Data["errmsg"] = "注册失败"
		beego.Error("插入失败")
		this.TplName = "register"
		return
	}

	//返回数据
	//this.Ctx.WriteString("注册成功")
	//注册之后跳转到登陆页面
	//这个方式只是渲染了页面但是没有跳转到登陆页面
	//this.TplName="login.html"
	//这个是跳转到登陆页面
	this.Redirect("/login", 302)
}

//展示登陆页面
func (this *UserController) ShowLogin() {
	dec := this.Ctx.GetCookie("userName")
	userName, _ := base64.StdEncoding.DecodeString(dec)
	if string(userName) != "" {
		this.Data["userName"] = string(userName)
		this.Data["checked"] = "checked"
	} else {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}
	this.TplName = "login.html"
}

//处理登陆业务
func (this *UserController) HandleLogin() {
	//1 接受数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	//2 校验数据
	if userName == "" || pwd == "" {
		beego.Error("登陆账号和密码为空")
		this.TplName = "login.html"
		this.Data["errmsg"] = "密码和账号错误"
		return
	}

	//3 处理数据

	//查询对象
	o := orm.NewOrm()
	//获取查询对象
	var user models.User
	//给查询对象赋值
	user.UserName = userName

	//查询
	//判断用户名是否正确
	err := o.Read(&user, "userName")
	if err != nil {
		this.Data["errmsg"] = "用户名不存在"
		this.TplName = "login.html"
		return
	}

	//判断密码是否正确

	if user.Pwd != pwd {
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}

	//获取是否记住用户名
	remember := this.GetString("remember")
	if remember == "on" {
		enc := base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName", enc, 3600*1)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}

	//返回数据
	//this.Ctx.WriteString("登陆成功")

	//this.TplName="login.html"
	this.SetSession("userName", userName)
	this.Redirect("/articleList", 302)

}
func (this *UserController)Logout(){
	fmt.Println("gfseawfgsreawrgseawrs f")
	//删除session
	this.DelSession("userName")
	fmt.Println("gfseawfgsreawrgseawrs f")

	this.Redirect("/login",302)
}