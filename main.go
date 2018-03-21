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
	beego.Router("/", &c.HomeController{})
	beego.Run()
}