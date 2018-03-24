package controllers

import (
	"roomreserve/helpers"

	"github.com/astaxie/beego"
)

//CalenendarController _
type CalenendarController struct {
	beego.Controller
}

//Get Home
func (c *CalenendarController) Get() {
	c.Data["title"] = "ปฏิทินการใช้งาน"
	c.Data["calendar"] = "active"
	c.Data["username"] = helpers.GetUser(c.Ctx.Request)
	c.Data["userimg"] = helpers.GetUserImage(c.Ctx.Request)
	c.Layout = "layout.html"
	c.TplName = "calendar/calendar.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "calendar/calendar-js.html"
	c.Render()
}
