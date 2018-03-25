package controllers

import (
	"html/template"
	"os"
	"path/filepath"
	"roomreserve/helpers"
	"roomreserve/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/google/uuid"
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
	file, err := models.GetReserveFile(int(ID))
	if err == nil {
		c.Ctx.Output.Download(file.FilePath1, file.FileName)
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
