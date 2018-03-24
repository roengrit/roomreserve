package main

import (
	"fmt"
	c "roomreserve/controllers"
	_ "roomreserve/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres",
		"host="+beego.AppConfig.String("pgurls")+
			" port="+beego.AppConfig.String("pgport")+
			" user="+beego.AppConfig.String("pguser")+
			" password="+beego.AppConfig.String("pgpass")+
			" dbname="+beego.AppConfig.String("pgdb")+
			" sslmode="+beego.AppConfig.String("pgsslmode"))
}

func main() {

	name := "default"
	force := false                             // Drop table and re-create.
	verbose := true                            // Print log.
	err := orm.RunSyncdb(name, force, verbose) // Error.
	if err != nil {
		fmt.Println(err)
	}
	beego.Router("/service/secure/json/", &c.ServiceController{}, "get:GetXSRF")
	beego.Router("/login", &c.UserController{})
	beego.Router("/logout", &c.LogoutController{})
	beego.Router("/forget-password", &c.ForgetController{})
	beego.Router("/", &c.HomeController{})
	beego.Router("/calendar", &c.CalenendarController{})

	beego.Router("/room/?:id", &c.RoomController{}, "get:Get;post:Post;delete:Delete")
	beego.Router("/room/read/?:id", &c.RoomReadController{}, "get:Read")
	beego.Router("/room/list", &c.RoomController{}, "get:RoomList;post:GetRoomListJSON")

	beego.Router("/user/?:id", &c.UserManageController{}, "get:Get;post:Post;delete:Delete")
	beego.Router("/profile", &c.UserManageController{}, "get:ProfileGet;post:ProfilePost")
	beego.Router("/user/list", &c.UserManageController{}, "get:UserList;post:GetUserListJSON")

	beego.Router("/role/?:id", &c.RoleController{}, "get:Get;post:Post;delete:Delete")
	beego.Router("/role/list", &c.RoleController{}, "get:RoleList;post:GetRoleListJSON")

	beego.Run()
}
