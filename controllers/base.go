package controllers

import (

	"github.com/astaxie/beego"
)
//BaseController _
type BaseController struct {
	beego.Controller
}

//Prepare _
func (b *BaseController) Prepare() {
}
