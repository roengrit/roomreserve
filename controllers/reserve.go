package controllers

import (
	"html/template"
	"roomreserve/models"

	"github.com/astaxie/beego"
)

//ReserveController _
type ReserveController struct {
	beego.Controller
}

//Get -
func (c *ReserveController) Get() {
	//roleID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	c.Data["title"] = "รายละเอียดการจอง"
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "reserve/reserve.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "reserve/reserve-js.html"
	c.Render()
}
