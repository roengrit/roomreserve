package controllers

import (
	"roomreserve/helpers"
	"roomreserve/models"
	"time"

	"github.com/astaxie/beego"
)

//HomeController _
type HomeController struct {
	beego.Controller
}

//Get Home
func (c *HomeController) Get() {
	c.Data["title"] = "จองห้องประชุม"
	c.Data["home"] = "active"
	c.Data["username"] = helpers.GetUser(c.Ctx.Request)
	c.Data["userimg"] = helpers.GetUserImage(c.Ctx.Request)
	c.Data["err"] = c.Ctx.Request.URL.Query().Get("err")
	c.Data["room"] = models.GetAllActiveRoom()
	c.Data["currentDate"] = time.Now().AddDate(543, 0, 0)
	c.Layout = "layout.html"
	c.TplName = "home/index.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "home/index-js.html"
	c.Render()
}
