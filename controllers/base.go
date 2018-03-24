package controllers

import (
	h "roomreserve/helpers"
	"roomreserve/models"
	"strings"

	"github.com/astaxie/beego"
)

//BaseController _
type BaseController struct {
	beego.Controller
}

//Prepare _
func (b *BaseController) Prepare() {
	val := h.GetUser(b.Ctx.Request)
	if val == "" {
		b.Ctx.Redirect(302, "/login?ref="+b.Ctx.Request.URL.RequestURI())
	}
	b.Data["username"] = h.GetUser(b.Ctx.Request)
	b.Data["userimg"] = h.GetUserImage(b.Ctx.Request)

	if val != "" {
		if strings.Contains(b.Ctx.Request.URL.Path, "user") {
			user, _ := models.GetUserByUserName(h.GetUser(b.Ctx.Request))
			s := strings.Split(user.Role.Access, ",")
			user.Role.User, user.Role.Role, user.Role.Room = s[0], s[1], s[2]
			if user.Role.User == "" {
				b.Ctx.Redirect(302, "/?err=คุณไม่มีสิทธิ์ใช้งาน")
			}
		}
		if strings.Contains(b.Ctx.Request.URL.Path, "role") {
			user, _ := models.GetUserByUserName(h.GetUser(b.Ctx.Request))
			s := strings.Split(user.Role.Access, ",")
			user.Role.User, user.Role.Role, user.Role.Room = s[0], s[1], s[2]
			if user.Role.Role == "" {
				b.Ctx.Redirect(302, "/?err=คุณไม่มีสิทธิ์ใช้งาน")
			}
		}
		if strings.Contains(b.Ctx.Request.URL.Path, "room") {
			user, _ := models.GetUserByUserName(h.GetUser(b.Ctx.Request))
			s := strings.Split(user.Role.Access, ",")
			user.Role.User, user.Role.Role, user.Role.Room = s[0], s[1], s[2]
			if user.Role.Room == "" {
				b.Ctx.Redirect(302, "/?err=คุณไม่มีสิทธิ์ใช้งาน")
			}
		}
	}

}
