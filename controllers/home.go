package controllers

import "github.com/astaxie/beego"

//HomeController _
type HomeController struct {
	beego.Controller
}

//Get Home
func (c *HomeController) Get() {
	c.Data["title"] = "หน้าหลัก"
	c.Layout = "layout.html"
	c.TplName = "home/index.html"
	c.Render()
}
