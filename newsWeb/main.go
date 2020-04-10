package main

import (
	"github.com/astaxie/beego"
	_ "newsWeb/models"
	_ "newsWeb/routers"
)

func main() {
	 // 在beego run 里面把函数关联起来
	beego.AddFuncMap("PrePage",PrePageIndex)
	beego.AddFuncMap("NextPage",NextPage)

	beego.Run()
}

//第二部 在代码里面定义一个函数
func PrePageIndex (pageIndex int) int{
	prePage:=pageIndex-1
	if prePage<1{
		prePage=1
	}
	return prePage

}

func NextPage(pageIndex int,pageCount float64)int  {

    if pageIndex+1>int(pageCount){
		return pageIndex
	}
	return pageIndex+1
}

//func fire (){
//	fmt.Println("d")
//}
//
//func f(){
//	//定义一个方法变量
//	var l func()
//
//	l=fire
//	//实现
//	l()
//}
//
