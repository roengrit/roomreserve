package controllers

import (
	"roomreserve/helpers"
)

//RoomController _
type RoomController struct {
	BaseController
}

//Get Home
func (c *RoomController) Get() {
	c.Data["title"] = "จองห้องประชุม"
	c.Data["username"] = helpers.GetUser(c.Ctx.Request)
	c.Layout = "layout.html"
	c.TplName = "room/room.html"
	//c.LayoutSections = make(map[string]string)
	//c.LayoutSections["scripts"] = "home/index-js.html"
	c.Render()
}
