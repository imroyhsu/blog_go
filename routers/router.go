package routers

import (
	"blog_go/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {

	// 这段代码放在router.go文件的init()的开头,解决前后端分离跨域问题
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: 	  []string{"http://"+beego.AppConfig.String("front_end_domain")+":"+beego.AppConfig.String("front_end_port")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	/**
	登录的
	*/
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})

	/**
	  文章的
	*/
	beego.Router("/articles", &controllers.ArticleController{})
	beego.Router("/articles/details", &controllers.ArticleController{}, "get:GetDetails")

	//官网rest风格示例
	//beego.Router("/api/list",&RestController{},"*:ListFood")
	//beego.Router("/api/create",&RestController{},"post:CreateFood")
	//beego.Router("/api/update",&RestController{},"put:UpdateFood")
	//beego.Router("/api/delete",&RestController{},"delete:DeleteFood")
}
