package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleReg")
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/articleList",&controllers.ArticleController{},"get:ShowArticleList;post:HandleSelect")
	beego.Router("/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandeAddArticle")
    beego.Router("/articleDetail",&controllers.ArticleController{},"get:ShowArticleDetail")
    beego.Router("/updateArticle",&controllers.ArticleController{},"get:ShowUpdateArticle;post:HandeUpdateArticle")
	beego.Router("/deleteArticle",&controllers.ArticleController{},"get:DeleteArticle")
	beego.Router("/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandeAddType")
	beego.Router("/deleteType",&controllers.ArticleController{},"get:DeleteType")
	beego.Router("/logout",&controllers.UserController{},"get:Logout")

}
