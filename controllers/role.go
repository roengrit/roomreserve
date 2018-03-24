package controllers

import (
	"bytes"
	"html/template"
	h "roomreserve/helpers"
	"roomreserve/models"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/form"
)

//RoleController _
type RoleController struct {
	BaseController
}

//RoleList _
func (c *RoleController) RoleList() {
	pn := h.NewPaging(0, 10, 0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["Page"] = 0
	c.Data["PerPage"] = 5
	c.Data["Pages"] = pn
	c.Data["title"] = "จัดการสิทธิ์"
	c.Layout = "layout.html"
	c.TplName = "role/role.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "role/role-js.html"
	c.Render()
}

//GetRoleListJSON _
func (c *RoleController) GetRoleListJSON() {
	searchTxt := c.Ctx.Request.FormValue("SearchTxt")
	page, _ := strconv.ParseInt(c.Ctx.Request.FormValue("Page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.Ctx.Request.FormValue("PerPage"), 10, 64)
	page = h.PrePaging(page)
	offset := h.CalOffsetPaging(page, perPage)
	ret := models.RoleListJSON{}
	num, list, _ := models.GetRoleList(uint(offset), uint(perPage), searchTxt)

	pn := h.NewPaging(page, perPage, num)
	ret.RoleList = &list

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
func (c *RoleController) Get() {
	roleID, _ := strconv.ParseInt(c.Ctx.Request.URL.Query().Get("id"), 10, 32)
	if roleID == 0 {
		c.Data["title"] = "สร้างสิทธิ์การใช้งาน"
	} else {
		c.Data["title"] = "แก้ไขสิทธิ์การใช้งาน"
		role, _ := models.GetRole(int(roleID))
		if role.Access != "" {
			s := strings.Split(role.Access, ",")
			role.User, role.Role, role.Room = s[0], s[1], s[2]
		}
		c.Data["m"] = role
	}
	c.Data["ret"] = models.RetModel{}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "role/create.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "role/create-js.html"
	c.Render()
}

//Post _
func (c *RoleController) Post() {
	var role models.Role
	decoder := form.NewDecoder()
	err := decoder.Decode(&role, c.Ctx.Request.Form)
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
	if ret.RetOK {
		role.Access = role.User + "," + role.Role + "," + role.Room
	}
	if ret.RetOK && role.ID == 0 {
		role.CreatedAt = time.Now()
		role.Creator = &actionUser
		id, err := models.CreateRole(role)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
			role.ID = int(id)
		}
	} else if ret.RetOK && role.ID > 0 {
		role.EditedAt = time.Now()
		role.Editor = &actionUser
		err := models.UpdateRole(role)
		if err != nil {
			ret.RetOK = false
			ret.RetData = err.Error()
		} else {
			ret.RetData = "บันทึกสำเร็จ"
		}
	}
	if role.ID == 0 {
		c.Data["title"] = "สร้างสิทธิ์การใช้งาน"
		c.Data["m"] = role
	} else {
		c.Data["title"] = "แก้ไขสิทธิ์การใช้งาน"
		role, _ := models.GetRole(int(role.ID))
		if role.Access != "" {
			s := strings.Split(role.Access, ",")
			role.User, role.Role, role.Room = s[0], s[1], s[2]
		}
		c.Data["m"] = role
	}
	c.Data["Role"] = models.GetAllRole()
	c.Data["ret"] = ret
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout.html"
	c.TplName = "role/create.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["scripts"] = "role/create-js.html"
	c.Render()
}

//Delete Delete
func (c *RoleController) Delete() {
	ID, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 32)
	ret := models.RetModel{}
	ret.RetOK = true
	err := models.DeleteRole(int(ID))
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
