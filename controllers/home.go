package controllers

import (
	"bytes"
	"html/template"
	"roomreserve/helpers"
	"roomreserve/models"
	"strconv"
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
	c.Data["room"] = models.GetAllRoom()
	c.Data["currentDate"] = time.Now().AddDate(543, 0, 0)
	c.Layout = "layout.html"
	c.TplName = "home/index.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "home/index-js.html"
	c.Render()
}

//Post Home
func (c *HomeController) Post() {
	searchTxt := c.Ctx.Request.FormValue("Title")
	room := c.Ctx.Request.FormValue("Room")
	status := c.Ctx.Request.FormValue("Status")

	dateBegin := c.Ctx.Request.FormValue("dateBeginPost")
	dateEnd := c.Ctx.Request.FormValue("dateEndPost")

	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = helpers.PrePaging(page)
	offset := helpers.CalOffsetPaging(page, perPage)

	ret := models.RoomReserveListJSON{}
	num, list, _ := models.GetReserveList(uint(offset), uint(perPage), searchTxt, status, room, dateBegin, dateEnd)

	pn := helpers.NewPaging(page, perPage, num)
	ret.RoomReserveList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}
