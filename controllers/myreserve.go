package controllers

import (
	"bytes"
	"html/template"
	"roomreserve/helpers"
	"roomreserve/models"
	"strconv"
	"time"
)

//MyReserveController _
type MyReserveController struct {
	BaseController
}

//Get Home
func (c *MyReserveController) Get() {
	c.Data["title"] = "รายการจองของฉััน"
	c.Data["reserve_list"] = "active"
	c.Data["err"] = c.Ctx.Request.URL.Query().Get("err")
	c.Data["room"] = models.GetAllRoom()
	c.Data["currentDate"] = time.Now().AddDate(543, 0, 0)
	c.Layout = "layout.html"
	c.TplName = "myreserve/index.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "myreserve/index-js.html"
	c.Render()
}

//Post Home
func (c *MyReserveController) Post() {
	searchTxt := c.Ctx.Request.FormValue("Title")
	room := c.Ctx.Request.FormValue("Room")
	status := c.Ctx.Request.FormValue("Status")

	dateBegin := c.Ctx.Request.FormValue("dateBeginPost")
	dateEnd := c.Ctx.Request.FormValue("dateEndPost")

	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = helpers.PrePaging(page)
	offset := helpers.CalOffsetPaging(page, perPage)

	ret := models.MyReserveListJSON{}
	num, list, _ := models.GetMyReserveList(uint(offset), uint(perPage), searchTxt, status, room, dateBegin, dateEnd)

	pn := helpers.NewPaging(page, perPage, num)
	ret.MyReserveList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}
