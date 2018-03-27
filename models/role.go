package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

//Role Role
type Role struct {
	ID        int
	Lock      bool
	Name      string    `orm:"size(225)"`
	Access    string    `orm:"size(500)"`
	User      string    `orm:"-"`
	Role      string    `orm:"-"`
	Room      string    `orm:"-"`
	HideTitle string    `orm:"-"`
	Creator   *User     `orm:"rel(fk)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	Editor    *User     `orm:"null;rel(fk)"`
	EditedAt  time.Time `orm:"null;auto_now;type(datetime)"`
}

//RoleListJSON -
type RoleListJSON struct {
	RoleList *[]Role
	Paging   string
	Page     uint
	PerPage  uint
}

func init() {
	orm.RegisterModel(new(Role))
}

//GetAllRole -
func GetAllRole() *[]Role {
	role := &[]Role{}
	o := orm.NewOrm()
	o.QueryTable("role").RelatedSel().All(role)
	return role
}

//GetRoleList `select * from x offset $1 limit $2`
func GetRoleList(currentPage, lineSize uint, term string) (num int64, roleListJSON []Role, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.name,
					T0.lock 					 
			   FROM "role" T0 
			   WHERE (lower(T0.name) like lower(?)) order by T0.name`
	num, _ = o.Raw(sql, "%"+term+"%").QueryRows(&roleListJSON)
	if lineSize+currentPage > uint(num) {
		lineSize = uint(num)
	} else if currentPage > 0 {
		lineSize = lineSize + currentPage
	}
	if currentPage > lineSize {
		currentPage = 0
	}
	roleListJSON = roleListJSON[currentPage:lineSize]
	return num, roleListJSON, err
}

//GetRole GetRole
func GetRole(ID int) (role *Role, errRet error) {
	Role := &Role{}
	o := orm.NewOrm()
	o.QueryTable("role").Filter("ID", ID).RelatedSel().One(Role)
	return Role, errRet
}

//CreateRole _
func CreateRole(Role Role) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&Role)
	o.Commit()
	return ID, err
}

//UpdateRole _
func UpdateRole(role Role) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetRole(role.ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if getUpdate == nil {
		errRet = errors.New("ไม่พบข้อมูล")
	} else if errRet == nil {
		role.CreatedAt = getUpdate.CreatedAt
		role.Creator = getUpdate.Creator
		if num, errUpdate := o.Update(&role); errUpdate != nil {
			errRet = errUpdate
			_ = num
		}
	}
	return errRet
}

//DeleteRole DeleteRole
func DeleteRole(ID int) (errRet error) {
	o := orm.NewOrm()
	unitDelete, _ := GetRoom(ID)
	if unitDelete.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if num, errDelete := o.Delete(&Role{ID: ID}); errDelete != nil && errRet == nil {
		errRet = errDelete
		_ = num
	}
	return errRet
}
