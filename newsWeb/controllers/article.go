package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"newsWeb/models"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

//展示/获取文章列表页面
func (this *ArticleController) ShowArticleList() {
	userName:=this.GetSession("userName")
	if userName==nil {
	this.Redirect("/login",302)
		return
	}
	this.Data["userName"]=userName.(string)

	// 查询数据库 拿出数据 传递视图
	//获取orm对象
	o := orm.NewOrm()
	//获取查询对象
	//var article [] models.Article
	article := make([]models.Article, 10)
	//查询
	//queryseter 座高级查询使用的数据类型
	//查询那张表
	qs := o.QueryTable("Article")

	//查询所有文章
	qs.All(&article)

	//beego.Info(article)

	this.Data["article"] = article

	//实现分页
	// 获取宗页数
	count, _ := qs.Count()
	//设置多少条数据一页
	pageSize := int64(2)

	//使用float64 类型
	pageCount := float64(count) / float64(pageSize)
	//获取总记录和宗页数
	//向上 去整数
	pageCount = math.Ceil(pageCount)
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	//获取首页和末叶数据
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		//this.Data["errmag"]="获取 分叶失败"
		//this.TplName="add.html"
		//return
		pageIndex = 1
	}

	beego.Info(pageIndex)
	//获取分页的数据
	start := pageSize * (int64(pageIndex) - 1)

	//数据库 部分查询
	 //orm是惰性查询方式不特别的指定就不会查询
	//RelatedSel 就是一对多关系表查询用来查询另外一张表
	qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&article)
	this.Data["pageIndex"] = pageIndex
	this.Data["article"] = article

	//获取数据
	var articleTypes [] models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
    this.Data["articleTypes"]=articleTypes
	this.TplName = "index.html"

}

//展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	//获取类型
	o:=orm.NewOrm()
	var articleTypes [] models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"]=articleTypes
	this.TplName = "add.html"
}

//处理添加文章业务
func (this *ArticleController) HandeAddArticle() {
	//接受数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	//校验数据
	if articleName == "" || content == "" {
		this.Data["errmsg"] = "文章标题或内容不呢个为空"
		this.TplName = "add.html"
		return
	}
     typeName:=this.GetString("select")
	//一接受图片
	file, head, err := this.GetFile("uploadname")
	if err != nil {
		this.Data["errmsg"] = "获取文件失败"
		this.TplName = "add.html"
		return
	}
	defer file.Close()

	//1 接受图片之前的操作 判断图片的大小
	if head.Size > 500000 {
		this.Data["errmsg"] = "文件太大，上传失败"
		this.TplName = "add.html"
		return
	}
	//2 判断图的格式
	//拿到文件的后缀
	fileExt := path.Ext(head.Filename)
	if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误"
		this.TplName = "add.html"
		return
	}
	//3 文件重复防止
	//"2006-01-02-15-04-05"这个时间是固定的
	fileName := time.Now().Format("2006-01-02-15-04-05") + fileExt
	// 参数1 文件的名字 参数2 文件的所要保存的路径  在这maodou里要加 . 表示但前路径
	this.SaveToFile("uploadname", "./static/image/"+fileName)
	//处理数据
	//插入 数据
	o := orm.NewOrm()
	var article models.Article
	//给对象赋值
	article.Title = articleName
	article.Content = content
	//在插入的时候我们不需要加.  "./static/image/"+fileName
	article.Image = "/static/image/" + fileName
	//根据类型名称获取类型对象
    var articleType models.ArticleType
	articleType.TypeName=typeName
	o.Read(&articleType,"TypeName")
	article.ArticleType=&articleType

	_, err = o.Insert(&article)
	if err != nil {

		this.Data["errmsg"] = "添加文章失败，请重新添加"
		this.TplName = "add.html"
		return
	}
	//返回页面

	this.Redirect("/articleList", 302)
}

func UpdateFile(this *ArticleController, filePath string) string {
	//一接受图片
	file, head, err := this.GetFile(filePath)
	if err != nil {
		this.Data["errmsg"] = "获取文件失败"
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()

	//1 接受图片之前的操作 判断图片的大小
	if head.Size > 500000 {
		this.Data["errmsg"] = "文件太大，上传失败"
		this.TplName = "add.html"
		return ""
	}
	//2 判断图的格式
	//拿到文件的后缀
	fileExt := path.Ext(head.Filename)
	if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误"
		this.TplName = "add.html"
		return ""
	}
	//3 文件重复防止
	//"2006-01-02-15-04-05"这个时间是固定的
	fileName := time.Now().Format("2006-01-02-15-04-05") + fileExt
	// 参数1 文件的名字 参数2 文件的所要保存的路径  在这maodou里要加 . 表示但前路径
	this.SaveToFile(filePath, "./static/image/"+fileName)
	return "/static/image/" + fileName

}
//展示文章详情页
func (this *ArticleController) ShowArticleDetail() {
	//获取数据
	articleId, err := this.GetInt("id")
	if err != nil {
		this.Data["errmsg"] = "请求路径错误"
		this.TplName = "index.html"
		return
	}
	//校验数据
	//处理数据
	//查询数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	err = o.Read(&article)
	if err != nil {
		this.Data["errmsg"] = "请求路径错误"
		this.TplName = "index.html"
		return
	}
   fmt.Println("gggggggggggggggggggggggg")
	//获取article对象
	//获取多对多操作对象
	m2m:=o.QueryM2M(&article,"User")
	//获取要插入的数据
	var user models.User
	userName:=this.GetSession("userName")
	//将借口类型转换为字符串类型
	user.UserName=userName.(string)
	o.Read(&user,"UserName")
	//插入多对多关系
	m2m.Add(user)
	//加载关系
	o.LoadRelated(&article,"User")
	//filter 过滤器 制定查询查询条件进行过滤
	var users []models.User
	//这里是双下划线
    o.QueryTable("User").Filter("Articles__Article__Id",articleId).Distinct().All(&users)
	this.Data["user"]=users
	//返回数据
	this.Data["article"] = article
	this.TplName = "content.html"

}

//展示编辑文章页面
func (this *ArticleController) ShowUpdateArticle() {
	//获取数据
	articleId, err := this.GetInt("id")

	errmsg := this.GetString("errmsg")
	if errmsg != "" {
		this.Data["errmsg"] = errmsg
	}
	//校验数据
	if err != nil {
		beego.Error("请求路径错误")
		this.Redirect("/articleList?errmsg", 302)
		return
	}
	beego.Info(articleId)

	//数据处理
	//查询操作
	o := orm.NewOrm()
	//获取orm对象
	var article models.Article
	article.Id = articleId
	o.Read(&article)
	//返回数据
	this.Data["article"] = article
	this.TplName = "update.html"
}

//处理编辑更新业务

func (this *ArticleController) HandeUpdateArticle() {

	//接受数据
	//接受文章的名字
	articleName := this.GetString("articleName")
	//接受文章的内容
	content := this.GetString("content")
	//接受文章的图片
	fileName := UpdateFile(this, "uploadname")

	articleId, err2 := this.GetInt("id")

	//校验数据
	if articleName == "" || content == "" || fileName == "" || err2 != nil {
		errmsg := "内容不能为空"
		this.Redirect("/updateArticle?id="+strconv.Itoa(articleId)+"&errmsg="+errmsg, 302)
		return

	}
	//处理数据
	//update 更新操作
	o := orm.NewOrm()
	var article models.Article

	article.Id = articleId
	err := o.Read(&article)
	if err != nil {
		errmsg := "更新文章为空"
		this.Redirect("/updateArticle?id="+strconv.Itoa(articleId)+"&errmsg="+errmsg, 302)
		return
	}
	//给更新对象赋值
	article.Title = articleName
	article.Content = content
	article.Image = fileName

	//更新
	o.Update(&article)
	//返回数据
	this.Redirect("/articleList", 302)

}

//删除业务处理
func (this *ArticleController) DeleteArticle() {

	//接受数据
	articleId, err := this.GetInt("id")
	//校验数据

	if err != nil {
		beego.Error("路径不对")
		this.Redirect("/articleList", 302)
		return
	}
	//处理了数据

	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	_, err = o.Delete(&article)
	if err != nil {
		beego.Error("del 失败")
		this.Redirect("/articleList", 302)
		return
	}

	//返回数据
	this.Redirect("/articleList", 302)

}

//展示添加分类
func (this *ArticleController)ShowAddType(){

	//接受数据
	o:=orm.NewOrm()
	//查询容器
	var articleTypes [] models.ArticleType
	//指定查询表
	qs:=o.QueryTable("ArticleType")
	qs.All(&articleTypes)
	 //返回给视图
	 this.Data["articleTypes"]=articleTypes
	 this.TplName="addType.html"
}

//处理类型添加业务
func  (this *ArticleController)HandeAddType(){
	//接受数据
	typeName:= this.GetString("typeName")
	//校验数据
	if typeName==""{
		this.Data["errmsg"]="分类为空"
		this.Redirect("/addType",302)
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var articleType  models.ArticleType
      articleType.TypeName=typeName
	_,err:=o.Insert(&articleType)
	if err!=nil {
		this.Data["errmsg"]="分类为空"
		//this.Redirect("/addType",302)
		this.TplName="addType.html"
		return
	}
	//返回数据
	this.Redirect("/addType",302)

}

//删除类型
func (this *ArticleController)DeleteType(){
    //接受数据
	 typeId,err:= this.GetInt("Id")


	//检验数据
	if err!=nil{
		beego.Error("删除失败")
		this.Redirect("/addhtml",302)
	}
	//处理数据
	o:=orm.NewOrm()
	//获取删除对象
    var articleType models.ArticleType
	//给产出对象赋值
	articleType.Id=typeId
	
	 _,err=o.Delete("")
	if  err!=nil {
		beego.Error("删除失败")
		this.Redirect("/addType",302)
		return
	}
	//返回数据,
	this.Redirect("/addType",302)

}

func  (this *ArticleController)HandleSelect()  {
	//接受数据
	typeName:= this.GetString("select")
	//校验数据

	//处理数据


	//返回数据


cp 

}