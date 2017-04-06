package routers

import (
	"travision/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/order", &controllers.OrderController{})
	beego.Router("/driver", &controllers.DriverController{})
	beego.Router("/carcoin", &controllers.CarcoinController{})
}
