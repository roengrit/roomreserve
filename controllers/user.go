package controllers

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	h "roomreserve/helpers"
	"roomreserve/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/go-playground/form"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//UserController _
type UserController struct {
	beego.Controller
}

//LogoutController _
type LogoutController struct {
	beego.Controller
}

//ForgetController _
type ForgetController struct {
	beego.Controller
}

//UserManageController _
type UserManageController struct {
	BaseController
}

//Get to view login
func (c *UserController) Get() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["ref"] = c.Ctx.Request.URL.Query().Get("ref")
	c.Data["username"] = "admin"
	c.Data["title"] = "เข้าสู่ระบบเพื่อเริ่มการทำงาน"
	c.TplName = "user/login.html"
	c.Render()
}

//Post to login
func (c *UserController) Post() {
	usernameForm := c.GetString("username")
	passwordForm := c.GetString("password")
	if ok, err := models.Login(usernameForm, passwordForm); ok {
		user, _ := models.GetUserByUserName(usernameForm)
		if ok, err = h.KeepLogin(c.Ctx.ResponseWriter, usernameForm, user.ID, user.Role.ID, user.ImagePath1); ok == true {
			if c.GetString("ref") != "" {
				c.Ctx.Redirect(http.StatusFound, c.GetString("ref"))
			}
			c.Ctx.Redirect(http.StatusFound, "/")
		} else {
			c.Data["error"] = err
		}
	} else {
		c.Data["error"] = err
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["username"] = c.GetString("username")
	c.Data["title"] = "เข้าสู่ระบบเพื่อเริ่มการทำงาน"
	c.TplName = "user/login.html"
	c.Render()
}

//ChangePass _
func (c *UserController) ChangePass() {
	val := h.GetUser(c.Ctx.Request)
	if val == "" {
		c.Ctx.Redirect(http.StatusFound, "/auth")
	}
	c.Data["UserDisplay"] = val
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["title"] = "เปลี่ยนรหัสผ่าน"
	c.Layout = "layout.html"
	c.TplName = "user/change-pass.html"
	c.Render()
}

//UpdatePass _
func (c *UserController) UpdatePass() {
	val := h.GetUser(c.Ctx.Request)
	if val == "" {
		c.Ctx.Redirect(http.StatusFound, "/login")
	}
	passwordForm := c.GetString("password")
	passwordReTryForm := c.GetString("password-retry")
	if passwordForm == passwordReTryForm && passwordForm != "" && len(passwordForm) >= 6 {
		val := h.GetUser(c.Ctx.Request)
		if ok, err := models.ChangePass(val, passwordForm); ok {
			c.Data["success"] = "ok"
		} else {
			c.Data["error"] = err
		}
	} else {
		c.Data["error"] = "รหัสผ่านไม่ตรงกัน"
	}
	if len(passwordForm) < 6 {
		c.Data["error"] = "รหัสผ่านต้องอย่างน้อย 6 ตัว"
	}
	if passwordForm == "" {
		c.Data["error"] = "กรุณาระบุรหัสผ่านให้ตรงกัน"
	}
	c.Data["UserDisplay"] = val
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["title"] = "เปลี่ยนรหัสผ่าน"
	c.Layout = "layout.html"
	c.TplName = "user/change-pass.html"
	c.Render()
}

//Get to logout
func (c *LogoutController) Get() {
	h.LogOut(c.Ctx.ResponseWriter)
	c.Ctx.Redirect(http.StatusFound, "/login")
}

//Get _
func (c *ForgetController) Get() {

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["username"] = ""
	c.Data["title"] = "กรอก email เพื่อรับรหัสผ่าน"
	c.TplName = "forget-password/forget.html"
	c.Render()
}

//Post to
func (c *ForgetController) Post() {
	usernameForm := c.GetString("username")
	newPass := models.RandStringRunes(8)
	if hasUser, errFindeuser := models.CheckUser(usernameForm); hasUser {
		if errSendMail := h.SendMail(usernameForm, newPass); errSendMail == "" {
			if ok, err := models.ForgetPass(usernameForm, newPass); ok {
				c.Data["success"] = "ส่งรหัสผ่านสำเร็จ"
			} else {
				c.Data["error"] = err
			}
		} else {
			c.Data["error"] = errSendMail
		}
	} else {
		c.Data["error"] = errFindeuser
	}

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["username"] = c.GetString("username")

	c.TplName = "forget-password/forget.html"
	c.Render()
}

//UserList _
func (c *UserManageController) UserList() {
	pn := h.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 5
	c.Data["Pages"] = pn
	c.Data["title"] = "จัดการผู้ใช้งาน"
	c.Layout = "layout.html"
	c.TplName = "user/user.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "user/user-js.html"
	c.Render()
}

//GetUserListJSON _
func (c *UserManageController) GetUserListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")
	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = h.PrePaging(page)
	offset := h.CalOffsetPaging(page, perPage)
	ret := models.UserListJSON{}
	num, list, _ := models.GetUserList(uint(offset), uint(perPage), searchTxt)

	pn := h.NewPaging(page, perPage, num)
	ret.UserList = &list

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	t, _ := template.ParseFiles("views/paging.html")
	var tpl bytes.Buffer
	t.Execute(&tpl, pn)
	ret.Paging = tpl.String()
	ret.Page = uint(num)
	c.Data["json"] = ret
	c.ServeJSON()
}

//Get -
func (c *UserManageController) Get() {
	userID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if userID == 0 {
		c.Data["title"] = "สร้างผู้ใช้งาน"
	} else {
		c.Data["title"] = "แก้ไขผู้ใช้งาน"
		user, _ := models.GetUserByID(int(userID))
		c.Data["m"] = user
	}

	c.Data["Role"] = models.GetAllRole()
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "user/create.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "user/create-js.html"
	c.Render()
}

//Post _
func (c *UserManageController) Post() {
	var user models.User
	decoder := form.NewDecoder()
	err := decoder.Decode(&user, c.Ctx.Request.Form)
	ret := models.RetModel{}
	actionUser, _ := models.GetUser(h.GetUser(c.Ctx.Request))
	ret.RetOK = true
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else if c.GetString("Name") == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อ"
	}

	if ret.RetOK && user.Username == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อผู้ใช้"
	}
	if ret.RetOK && user.Password != "" {
		if user.Password == user.RePassword {
			if hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
				ret.RetOK = false
				ret.RetData = err.Error()
			} else {
				user.Password = string(hash)
			}
		} else {
			ret.RetOK = false
			ret.RetData = "การยืนยันรหัสผ่านไม่ตรงกัน"
		}
	}

	if ret.RetOK {
		file, header, _ := c.GetFile("imagePath1")
		if file != nil {
			fileName := header.Filename
			fileName = uuid.New().String() + filepath.Ext(fileName)
			filePathSave := "static/image/user/" + fileName
			err = c.SaveToFile("imagePath1", filePathSave)
			filePathSave = "/" + filePathSave
			if err == nil {
				user.ImagePath1 = filePathSave
				user.DeleteImage1 = 0
			}
		}
		_ = file
	}

	if ret.RetOK && user.ID == 0 {
		user.CreatedAt = time.Now()
		user.Creator = actionUser.ID
		id, err := models.CreateUser(user)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
			user.ID = int(id)
		}
	} else if ret.RetOK && user.ID > 0 {
		user.EditedAt = time.Now()
		user.Editor = actionUser.ID
		err := models.UpdateUser(user)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	}
	if user.ID == 0 {
		c.Data["title"] = "สร้างผู้ใช้งาน"
		c.Data["m"] = user
	} else {
		c.Data["title"] = "แก้ไขผู้ใช้งาน"
		user, _ := models.GetUserByID(int(user.ID))
		c.Data["m"] = user
	}
	c.Data["Role"] = models.GetAllRole()
	c.Data["ret"] = ret
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "user/create.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "user/create-js.html"
	c.Render()
}

//ProfileGet -
func (c *UserManageController) ProfileGet() {
	c.Data["title"] = "แก้ไขโปรไฟล์"
	user, _ := models.GetUser(h.GetUser(c.Ctx.Request))
	c.Data["m"] = user
	c.Data["important"] = "readonly"
	c.Data["Role"] = models.GetAllRole()
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "user/create.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "user/create-js.html"
	c.Render()
}

//ProfilePost _
func (c *UserManageController) ProfilePost() {

	var user models.User
	decoder := form.NewDecoder()
	err := decoder.Decode(&user, c.Ctx.Request.Form)
	ret := models.RetModel{}
	actionUser, _ := models.GetUser(h.GetUser(c.Ctx.Request))

	ret.RetOK = true
	if err != nil {
		ret.RetOK = false
		ret.RetData = err.Error()
	} else if c.GetString("Name") == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อ"
	}

	if ret.RetOK && user.Username == "" {
		ret.RetOK = false
		ret.RetData = "กรุณาระบุชื่อผู้ใช้"
	}
	if ret.RetOK && user.Password != "" {
		if user.Password == user.RePassword {
			if hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
				ret.RetOK = false
				ret.RetData = err.Error()
			} else {
				user.Password = string(hash)
			}
		} else {
			ret.RetOK = false
			ret.RetData = "การยืนยันรหัสผ่านไม่ตรงกัน"
		}
	}

	if ret.RetOK {
		file, header, _ := c.GetFile("imagePath1")
		if file != nil {
			fileName := header.Filename
			fileName = uuid.New().String() + filepath.Ext(fileName)
			filePathSave := "static/image/user/" + fileName
			err = c.SaveToFile("imagePath1", filePathSave)
			filePathSave = "/" + filePathSave
			if err == nil {
				user.ImagePath1 = filePathSave
				user.DeleteImage1 = 0
			}
		}
		_ = file
	}
	user.ID = actionUser.ID
	user.Username = actionUser.Username
	user.Role = actionUser.Role
	user.Active = actionUser.Active
	if ret.RetOK && user.ID > 0 {
		user.EditedAt = time.Now()
		user.Editor = actionUser.ID
		err := models.UpdateUser(user)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
			c.Data["m"] = user
		} else {
			ret.RetData = "บันทึกสำเร็จ"
			user, _ := models.GetUserByID(int(user.ID))
			c.Data["m"] = user
		}
	}
	c.Data["title"] = "แก้ไขโปรไฟล์"
	c.Data["important"] = "disabled"
	c.Data["Role"] = models.GetAllRole()
	c.Data["ret"] = ret
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "user/create.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "user/create-js.html"
	c.Render()
}

//Delete Delete
func (c *UserManageController) Delete() {
	ID, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 32)
	ret := models.RetModel{}
	ret.RetOK = true
	err := models.DeleteUser(int(ID))
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
