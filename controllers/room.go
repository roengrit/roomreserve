package controllers

import (
	"bytes"
	"html/template"
	"path/filepath"
	"roomreserve/helpers"
	"roomreserve/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"

	"github.com/go-playground/form"
	"github.com/google/uuid"
)

//RoomController _
type RoomController struct {
	BaseController
}

//RoomReadController _
type RoomReadController struct {
	beego.Controller
}

//Get Home
func (c *RoomController) Get() {
	roomID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if roomID == 0 {
		c.Data["title"] = "สร้างห้อง"
	} else {
		c.Data["title"] = "แก้ไขห้อง"
		room, _ := models.GetRoom(int(roomID))
		c.Data["m"] = room
	}
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "room/create.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "room/create-js.html"
	c.Render()
}

//Get Home
func (c *RoomReadController) Read() {
	roomID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	room, _ := models.GetRoom(int(roomID))
	c.Data["m"] = room
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "room/view.html"
	c.LayoutSections = make(map[string]string)
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
	for index := 1; index <= 6; index++ {
		file, header, _ := c.GetFile("imagePath" + strconv.Itoa(index))
		if file != nil {
			fileName := header.Filename
			fileName = uuid.New().String() + filepath.Ext(fileName)
			filePathSave := "static/image/room/" + fileName
			err = c.SaveToFile("imagePath"+strconv.Itoa(index), filePathSave)
			filePathSave = "/" + filePathSave
			if err == nil {
				switch "imagePath" + strconv.Itoa(index) {
				case "imagePath1":
					room.ImagePath1 = filePathSave
					room.DeleteImage1 = 0
				case "imagePath2":
					room.ImagePath2 = filePathSave
					room.DeleteImage2 = 0
				case "imagePath3":
					room.ImagePath3 = filePathSave
					room.DeleteImage3 = 0
				case "imagePath4":
					room.ImagePath4 = filePathSave
					room.DeleteImage4 = 0
				case "imagePath5":
					room.ImagePath5 = filePathSave
					room.DeleteImage5 = 0
				case "imagePath6":
					room.ImagePath6 = filePathSave
					room.DeleteImage6 = 0
				}
			}
		}
		_ = file
	}
	if ret.RetOK && room.ID == 0 {
		room.CreatedAt = time.Now()
		room.Creator = &actionUser
		id, err := models.CreateRoom(room)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
			room.ID = int(id)
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
	if room.ID == 0 {
		c.Data["title"] = "สร้างห้อง"
		c.Data["m"] = room
	} else {
		c.Data["title"] = "แก้ไขห้อง"
		room, _ := models.GetRoom(int(room.ID))
		c.Data["m"] = room
	}
	c.Data["ret"] = ret
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "room/create.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "room/create-js.html"
	c.Render()
}

//Delete Delete
func (c *RoomController) Delete() {
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
