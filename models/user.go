package models

import (
	"errors"
	"math/rand"
	"time"

	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
)

//User User
type User struct {
	ID           int
	Lock         bool
	Username     string `orm:"size(50)"`
	Password     string `orm:"size(255)"`
	RePassword   string `orm:"-"`
	Name         string `orm:"size(500)"`
	Depart       string `orm:"size(255)"`
	Tel          string `orm:"size(255)"`
	Email        string `orm:"size(255)"`
	Facebook     string `orm:"size(255)"`
	Line         string `orm:"size(255)"`
	ImagePath1   string `orm:"size(255)"`
	DeleteImage1 int    `orm:"-"`
	Role         *Role  `orm:"rel(fk)"`
	Active       bool
	CreatedAt    time.Time `orm:"auto_now_add;type(datetime)"`
	Creator      int
	EditedAt     time.Time `orm:"null;auto_now;type(datetime)"`
	Editor       int
}

//UserListJSON -
type UserListJSON struct {
	UserList *[]User
	Paging   string
	Page     uint
	PerPage  uint
}

func init() {
	orm.RegisterModel(new(User))
}

//Login Login
func Login(username, password string) (ok bool, errRet string) {
	o := orm.NewOrm()
	user := User{Username: username}
	err := o.Read(&user, "Username")
	if err == orm.ErrNoRows {
		errRet = "รหัสผู้ใช้/รหัสผ่านไม่ถูกต้อง"
	} else {
		if errCript := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); errCript != nil {
			errRet = "รหัสผู้ใช้/รหัสผ่านไม่ถูกต้อง"
		} else {
			ok = true
		}
	}
	return ok, errRet
}

//GetUser GetUser
func GetUser(username string) (userRet User, errRet error) {
	o := orm.NewOrm()
	user := User{Username: username}
	errRet = o.Read(&user, "Username")
	if errRet == orm.ErrNoRows {
		errRet = errors.New("ไม่พบผู้ใช้งานนี้ในระบบ")
	}
	return user, errRet
}

//GetUserList `select * from x offset $1 limit $2`
func GetUserList(currentPage, lineSize uint, term string) (num int64, userListJSON []User, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.name,
					T0.lock	,
					T0.username,
					T0.tel,
					T0.depart, 
					T0.active					 
			   FROM  "user" T0	 JOIN "role" T1 ON T0.role_id = T1.i_d
			   WHERE (lower(T0.name) like lower(?)
			   or  lower(T0.depart) like lower(?)  
			   ) order by T0.name`
	num, _ = o.Raw(sql, "%"+term+"%", "%"+term+"%").QueryRows(&userListJSON)
	if lineSize+currentPage > uint(num) {
		lineSize = uint(num)
	} else if currentPage > 0 {
		lineSize = lineSize + currentPage
	}
	if currentPage > lineSize {
		currentPage = 0
	}
	userListJSON = userListJSON[currentPage:lineSize]
	return num, userListJSON, err
}

//GetUserByID GetUserByID
func GetUserByID(ID int) (user *User, errRet error) {
	o := orm.NewOrm()
	userGet := &User{}
	o.QueryTable("user").Filter("ID", ID).RelatedSel().One(userGet)
	if nil != userGet {
		userGet.Password = ""
	} else {
		errRet = errors.New("ชื่อผู้ใช้หรือรหัสผ่านผิดพลาด")
	}
	return userGet, errRet
}

//GetUserPassword -
func GetUserPassword(ID int) (user *User, errRet error) {
	o := orm.NewOrm()
	userGet := &User{}
	o.QueryTable("user").Filter("ID", ID).RelatedSel().One(userGet)
	if nil != userGet {

	} else {
		errRet = errors.New("ชื่อผู้ใช้หรือรหัสผ่านผิดพลาด")
	}
	return userGet, errRet
}

//GetUserByUserName _
func GetUserByUserName(username string) (user *User, errRet error) {
	userGet := &User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("Username", username).RelatedSel().One(userGet)
	if nil != userGet {
		userGet.Password = ""
	} else {
		errRet = errors.New("ชื่อผู้ใช้หรือรหัสผ่านผิดพลาด")
	}
	return userGet, errRet
}

//ForgetPass ForgetPass
func ForgetPass(username, newpass string) (ok bool, errRet error) {
	o := orm.NewOrm()
	user := User{Username: username}
	errRet = o.Read(&user, "Username")
	if errRet == orm.ErrNoRows {
		errRet = errors.New("ชื่อผู้ใช้หรือรหัสผ่านผิดพลาด")
	} else {
		if hash, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost); err != nil {
			errRet = err
		} else {
			user.Password = string(hash)
			if _, errUpdate := o.Update(&user); errUpdate != nil {
				errRet = errUpdate
			} else {
				ok = true
			}
		}
	}
	return ok, errRet
}

//ChangePass ChangePass
func ChangePass(username, newpass string) (ok bool, errRet error) {
	o := orm.NewOrm()
	user := User{Username: username}
	errRet = o.Read(&user, "Username")
	if errRet == orm.ErrNoRows {
		errRet = errors.New("ชื่อผู้ใช้หรือรหัสผ่านผิดพลาด")
	} else {
		if hash, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost); err != nil {
			errRet = err
		} else {
			user.Password = string(hash)
			if num, errUpdate := o.Update(&user); errUpdate != nil {
				errRet = errUpdate
				_ = num
			} else {
				ok = true
			}
		}
	}
	return ok, errRet
}

//RandStringRunes _
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//CheckUser _
func CheckUser(username string) (ok bool, errRet error) {
	o := orm.NewOrm()
	user := User{Username: username}
	errRet = o.Read(&user, "Username")
	ok = true
	if errRet == orm.ErrNoRows {
		errRet = errors.New("ชื่อผู้ใช้หรือรหัสผ่านผิดพลาด")
		ok = false
	}
	if ok && user.Active == false {
		ok = false
		errRet = errors.New("ชื่อผู้ใช้หรือรหัสผ่านผิดพลาด")
	}
	return ok, errRet
}

//CreateUser _
func CreateUser(User User) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&User)
	o.Commit()
	return ID, err
}

//UpdateUser _
func UpdateUser(user User) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetUserPassword(user.ID)
	if user.DeleteImage1 == 1 {
		user.ImagePath1 = ""
	} else if user.ImagePath1 == "" {
		user.ImagePath1 = getUpdate.ImagePath1
	}
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if user.Password == "" {
		user.Password = getUpdate.Password
	}
	if getUpdate == nil {
		errRet = errors.New("ไม่พบข้อมูล")
	} else if errRet == nil {
		user.CreatedAt = getUpdate.CreatedAt
		if num, errUpdate := o.Update(&user); errUpdate != nil {
			errRet = errUpdate
			_ = num
		}
	}
	return errRet
}

//DeleteUser -
func DeleteUser(ID int) (errRet error) {
	o := orm.NewOrm()
	unitDelete, _ := GetUserByID(ID)
	if unitDelete.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if num, errDelete := o.Delete(&User{ID: ID}); errDelete != nil && errRet == nil {
		errRet = errDelete
		_ = num
	}
	return errRet
}
