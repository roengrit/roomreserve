package controllers

import "github.com/astaxie/beego"

//HomeController _
type HomeController struct {
	beego.Controller
}

//Get Home
func (c *HomeController) Get() {
	c.Data["title"] = "ค้นหาห้องว่าง"
	c.Layout = "layout.html"
	c.TplName = "home/index.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "home/index-js.html"
	c.Render()
}
