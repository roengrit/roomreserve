package controllers

import (
	"html/template"
	"os"
	"path/filepath"
	"roomreserve/helpers"
	"roomreserve/models"
	"strconv"
	"time"

	"github.com/go-playground/form"
	"github.com/google/uuid"
)

//ReserveController _
type ReserveController struct {
	BaseController
}

//Get -
func (c *ReserveController) Get() {
	ID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	c.Data["title"] = "รายละเอียดการจอง"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))
	c.Data["editUser"] = 0
	if ID != 0 {
		reserv, _ := models.GetReserveRoom(int(ID))
		c.Data["m"] = reserv
		if reserv.Creator.ID == actionUser.ID {
			c.Data["editUser"] = 1
		} else {
			c.Ctx.Redirect(302, "/?err=คุณไม่มีสิทธิ์ในการแก้ใขการจอง")
		}
	}
	c.Data["username"] = helpers.GetUser(c.Ctx.Request)
	c.Data["userimg"] = helpers.GetUserImage(c.Ctx.Request)
	c.Layout = "layout.html"
	c.TplName = "reserve/reserve.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "reserve/reserve-js.html"
	c.Render()
}

//Post -
func (c *ReserveController) Post() {
	var reserve models.RoomReserve
	decoder := form.NewDecoder()
	err := decoder.Decode(&reserve, c.Ctx.Request.Form)
	ret := models.RetModel{RetOK: true}
	actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))
	if actionUser.ID == 0 {
		ret.RetOK = false
		ret.RetData = "กรุณาเข้าสู่ระบบ"
	}
	if ret.RetOK && err == nil {
		reserve.DateBegin, err = helpers.CreateDateTimeFromString(c.GetString("DateBegin"))
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		}
	}

	if ret.RetOK && err == nil {
		reserve.DateEnd, err = helpers.CreateDateTimeFromString(c.GetString("DateEnd"))
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		}
	}

	if err == nil {
		if reserve.ID == 0 {
			reserve.CreatedAt = time.Now()
			reserve.Creator = &actionUser
			if ID, err := models.CreateReserveRoom(reserve); err == nil {
				reserve.ID = int(ID)
				ret.RetOK = true
				ret.RetData = "บันทึกสำเร็จ"
			} else {
				ret.RetOK = false
				ret.RetData = err.Error()
			}
		} else {
			reserve.EditedAt = time.Now()
			reserve.Editor = &actionUser

			if chk, chkErr := models.GetReserveRoom(reserve.ID); chkErr == nil {
				if chk.Creator.ID == actionUser.ID {
					if err := models.UpdateReserveRoom(reserve); err == nil {
						ret.RetOK = true
						ret.RetData = "บันทึกสำเร็จ"
					} else {
						ret.RetOK = false
						ret.RetData = err.Error()
					}
				} else {
					ret.RetOK = false
					ret.RetData = "คุณไม่มีสิทธิ์ในเอกสารฉบับนี้"
				}
			} else {
				ret.RetOK = false
				ret.RetData = "ไม่พบข้อมูล"
			}

		}
	} else {
		ret.RetOK = false
		ret.RetData = "เกิดข้อผิดพลาด"
	}
	if reserve.ID != 0 {
		reserv, _ := models.GetReserveRoom(int(reserve.ID))
		c.Data["m"] = reserv
	} else {
		c.Data["m"] = reserve
	}

	c.Data["ret"] = ret
	c.Data["title"] = "รายละเอียดการจอง"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "reserve/reserve.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "reserve/reserve-js.html"
	c.Render()
}

//FileAtt FileAtt
func (c *ReserveController) FileAtt() {
	file, header, _ := c.GetFile("file-att")
	ret := models.RetModel{}
	ret.RetOK = false
	if file != nil {
		fileName := header.Filename
		fileName = uuid.New().String() + filepath.Ext(fileName)
		filePathSave := "data/file/" + fileName
		err := c.SaveToFile("file-att", filePathSave)
		ret.RetData = filePathSave
		ret.Name = header.Filename
		if err != nil {
			ret.RetOK = false
			ret.Alert = err.Error()

		} else {
			ret.RetOK = true
			filedb := models.RoomReserveFile{}
			ID, _ := strconv.ParseInt(c.GetString("ID"), 10, 32)
			filedb.ReserveID = int(ID)
			filedb.FilePath1 = filePathSave
			filedb.FileName = header.Filename
			filedb.Status = 1
			filedb.CreatedAt = time.Now()
			actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))
			filedb.Creator = &actionUser
			ID, err = models.CreateReserveFile(filedb)
			if err != nil {
				ret.RetOK = false
				ret.Alert = err.Error()
			}
			ret.ID = ID
		}
	}
	_ = file
	c.Data["json"] = ret
	c.ServeJSON()
}

//FileDownload FileDownload
func (c *ReserveController) FileDownload() {
	id := c.Ctx.Input.Param(":id")
	ID, _ := strconv.ParseInt(id, 10, 32)
	actionUser, _ := models.GetUser(helpers.GetUser(c.Ctx.Request))
	file, err := models.GetReserveFile(int(ID))
	if err == nil {
		reserve, err := models.GetReserveRoom(file.ReserveID)
		if err == nil {
			if reserve.HideFile == 1 && actionUser.ID != 0 {
				c.Ctx.Output.Download(file.FilePath1, file.FileName)
			} else if reserve.HideFile == 0 {
				c.Ctx.Output.Download(file.FilePath1, file.FileName)
			}
		}
	}
}

//FileDelete -
func (c *ReserveController) FileDelete() {
	id := c.Ctx.Input.Param(":id")
	ID, _ := strconv.ParseInt(id, 10, 32)
	file, err := models.GetReserveFile(int(ID))
	errDelete := models.DeleteReserveFile(int(ID))
	ret := models.RetModel{}
	ret.RetOK = true
	if err == nil && errDelete == nil {
		err = os.Remove(file.FilePath1)
	} else {
		ret.RetOK = false
		ret.Alert = "ไม่สามารถลบไฟล์ได้"
	}
	c.Data["json"] = ret
	c.ServeJSON()
}
