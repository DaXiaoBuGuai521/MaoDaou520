package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	UserName string  `orm:"unique"`
	Pwd    string
	//设置多对多关系  用户和文章的关系
	//用户可以看多篇文章 一篇文章也可以被多个用户看
	Articles []*Article `orm:"rel(m2m)"`
}

type Article struct {
	//属性设置 orm后面就是属性设置  “pk;auto” Id  int `orm:""`
	Id int `orm:"pk;auto"`
	Title string  `orm:"size(100)"`
	Content string   `orm:"size(500)"`
	Time time.Time   `orm:"type(datetime);auto_now"`
	ReadCount int   `orm:"default(0)"`
	Image string  `orm:"null"`
	//外健设置
    ArticleType *ArticleType `orm:"rel(fk)"`
	User [] *User `orm:"reverse(many)"`
}

type ArticleType struct {
	Id int
	TypeName string  `orm:"size(100)"`
	//外见约束
	Articles []*Article `orm:"reverse(many)"`

}
func init(){


	//1 注册数据库
	//orm.RegisterDriver("mysql", orm.DRMySQL) //可以不加
	orm.RegisterDataBase("default", "mysql", "king:123456@tcp(192.168.31.165:3306)/test?charset=utf8")
	//ORM 必须注册一个别名为default的数据库，作为默认使用。
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//3 运行  //
	//orm.RunSyncdb("default", false, true)
	//修改参数2 为true 强制更新表
	//第二个参数为true的时候 强制更新 表的数据会消失
	orm.RunSyncdb("default", false, true)

}
