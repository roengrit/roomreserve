package controllers

import (
	"bytes"
	"html/template"
	"roomreserve/helpers"
	"roomreserve/models"
	"strconv"
	"time"

	"github.com/go-playground/form"
)

//RoomController _
type RoomController struct {
	BaseController
}

//Get Home
func (c *RoomController) Get() {
	cateID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if cateID == 0 {
		c.Data["title"] = "สร้างห้อง"
	} else {
		c.Data["title"] = "แก้ไขห้อง"
		category, _ := models.GetRoom(int(cateID))
		c.Data["m"] = category
	}
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "room/create.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "room/create-js.html"
	c.Render()
}

//RoomList _
func (c *RoomController) RoomList() {
	pn := helpers.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 5
	c.Data["Pages"] = pn
	c.Data["title"] = "จัดการห้อง"
	c.Layout = "layout.html"
	c.TplName = "room/room.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "room/room-js.html"
	c.Render()
}

//GetRoomListJSON _
func (c *RoomController) GetRoomListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")
	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = helpers.PrePaging(page)
	offset := helpers.CalOffsetPaging(page, perPage)
	ret := models.RoomListJSON{}
	num, list, _ := models.GetRoomList(uint(offset), uint(perPage), searchTxt)

	pn := helpers.NewPaging(page, perPage, num)
	ret.RoomList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//Post _
func (c *RoomController) Post() {

	var room models.Room
	decoder := form.NewDecoder()
	err := decoder.Decode(&room, c.Ctx.Request.Form)
	ret := models.RetModel{}
	actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))
	ret.RetOK = true
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else if c.GetString("Name") == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อ"
	}

	if ret.RetOK && room.ID == 0 {
		room.CreatedAt = time.Now()
		room.Creator = &actionUser
		_, err := models.CreateRoom(room)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	} else if ret.RetOK && room.ID > 0 {
		room.EditedAt = time.Now()
		room.Editor = &actionUser
		err := models.UpdateRoom(room)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	}
	ret.XSRF = c.XSRFToken()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["json"] = ret
	c.ServeJSON()
}

//DeleteRoom DeleteRoom
func (c *RoomController) DeleteRoom() {
	ID, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 32)
	ret := models.RetModel{}
	ret.RetOK = true
	err := models.DeleteRoom(int(ID))
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else {
		ret.RetData = "ลบข้อมูลสำเร็จ"
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["json"] = ret
	c.ServeJSON()
}
