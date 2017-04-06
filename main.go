package main

import (
	_ "travision/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/Dashboard", "static/Dashboard")
	beego.SetStaticPath("/image", "static/img")
	beego.Run()
}
