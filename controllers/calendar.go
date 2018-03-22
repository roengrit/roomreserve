package controllers

import "github.com/astaxie/beego"

//CalenendarController _
type CalenendarController struct {
	beego.Controller
}

//Get Home
func (c *CalenendarController) Get() {
	c.Data["title"] = "ปฏิทินการใช้งาน"
	c.Layout = "layout.html"
	c.TplName = "calendar/calendar.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "calendar/calendar-js.html"
	c.Render()
}
