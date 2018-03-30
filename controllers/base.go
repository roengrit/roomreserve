package controllers

import (
	"fmt"
	h "roomreserve/helpers"
	"roomreserve/models"
	"strings"

	"github.com/astaxie/beego"
)

//BaseController _
type BaseController struct {
	beego.Controller
}

//BaseNoAuthController _
type BaseNoAuthController struct {
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
		user, _ := models.GetUserByUserName(h.GetUser(b.Ctx.Request))
		s := strings.Split(user.Role.Access, ",")
		user.Role.User, user.Role.Role, user.Role.Room, user.Role.HideTitle = s[0], s[1], s[2], s[3]
		if strings.Contains(b.Ctx.Request.URL.Path, "user") {
			if user.Role.User == "" {
				b.Ctx.Redirect(302, "/?err=คุณไม่มีสิทธิ์ใช้งาน")
			}
		}
		if strings.Contains(b.Ctx.Request.URL.Path, "role") {
			if user.Role.Role == "" {
				b.Ctx.Redirect(302, "/?err=คุณไม่มีสิทธิ์ใช้งาน")
			}
		}
		if strings.Contains(b.Ctx.Request.URL.Path, "room") {
			if user.Role.Room == "" {
				b.Ctx.Redirect(302, "/?err=คุณไม่มีสิทธิ์ใช้งาน")
			}
		}
		if strings.Contains(b.Ctx.Request.URL.Path, "reserve") {
			b.Data["hideTitle"] = user.Role.HideTitle
		}
	}

}

//Prepare _
func (b *BaseNoAuthController) Prepare() {
	val := h.GetUser(b.Ctx.Request)
	b.Data["username"] = h.GetUser(b.Ctx.Request)
	b.Data["userimg"] = h.GetUserImage(b.Ctx.Request)
	fmt.Println(val)
	if val != "" {
		user, _ := models.GetUserByUserName(h.GetUser(b.Ctx.Request))
		s := strings.Split(user.Role.Access, ",")
		user.Role.User, user.Role.Role, user.Role.Room, user.Role.HideTitle = s[0], s[1], s[2], s[3]
		b.Data["hideTitle"] = user.Role.HideTitle
		b.Data["room_manage"] = user.Role.Room
		b.Data["role_manage"] = user.Role.Role
		b.Data["user_manage"] = user.Role.User
	}

}
