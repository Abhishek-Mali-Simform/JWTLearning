package routers

import (
	"JWTLearning/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/home", &controllers.MainController{}, "*:Home")
	beego.Router("/refresh", &controllers.MainController{}, "*:Refresh")
}
