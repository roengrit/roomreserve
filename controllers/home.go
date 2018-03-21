package controllers

//HomeController _
type HomeController struct {
	BaseController
}

//Get Home
func (c *HomeController) Get() {
	c.Data["title"] = "หน้าหลัก"
	c.Layout = "layout.html"
	c.TplName = "home/index.html"
	c.Render()
}
