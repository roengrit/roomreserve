package controllers

import (
	h "roomreserve/helpers"

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
		b.Ctx.Redirect(302, "/login?ref="+b.Ctx.Request.URL.Path)
	}
	b.Data["username"] = val
}
