package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

//Room _
type Room struct {
	ID              int
	Lock            bool
	Name            string `orm:"size(300)"`
	SupportText     string `orm:"size(300)"`
	LocationText    string `orm:"size(300)"`
	AddOnDeviceText string `orm:"size(300)"`
	RoomAdminText   string `orm:"size(300)"`
	ImagePath1      string `orm:"size(300)"`
	ImagePath2      string `orm:"size(300)"`
	ImagePath3      string `orm:"size(300)"`
	ImagePath4      string `orm:"size(300)"`
	ImagePath5      string `orm:"size(300)"`
	ImagePath6      string `orm:"size(300)"`
	DeleteImage1    bool   `orm:"-"`
	DeleteImage2    bool   `orm:"-"`
	DeleteImage3    bool   `orm:"-"`
	DeleteImage4    bool   `orm:"-"`
	DeleteImage5    bool   `orm:"-"`
	DeleteImage6    bool   `orm:"-"`
	Active          bool
	Creator         *User     `orm:"rel(fk)"`
	CreatedAt       time.Time `orm:"auto_now_add;type(datetime)"`
	Editor          *User     `orm:"null;rel(fk)"`
	EditedAt        time.Time `orm:"null;auto_now;type(datetime)"`
}

//RoomListJSON RoomListJSON
type RoomListJSON struct {
	RoomList *[]Room
	Paging   string
	Page     uint
	PerPage  uint
}

func init() {
	orm.RegisterModel(new(Room))
}

//GetAllRoom GetAllRoom
func GetAllRoom() *[]Room {
	room := &[]Room{}
	o := orm.NewOrm()
	o.QueryTable("room").RelatedSel().All(room)
	return room
}

//GetRoomList `select * from x offset $1 limit $2`
func GetRoomList(currentPage, lineSize uint, term string) (num int64, roomListJSON []Room, err error) {
	o := orm.NewOrm()
	var sql = `SELECT 
					T0.i_d,
					T0.name,
					T0.lock	,
					T0.support_text,
					T0.location_text,
					T0.add_on_device_text,
					T0.room_admin_text,
					T0.active	 				 ,
					T0.image_path1
			   FROM room T0	    
			   WHERE (lower(T0.name) like lower(?)) `

	num, _ = o.Raw(sql, "%"+term+"%").QueryRows(&roomListJSON)
	sql += " order by T0.name "

	if lineSize+currentPage > uint(num) {
		lineSize = uint(num)
	} else if currentPage > 0 {
		lineSize = lineSize + currentPage
	}
	if currentPage > lineSize {
		currentPage = 0
	}

	roomListJSON = roomListJSON[currentPage:lineSize]
	return num, roomListJSON, err
}

//GetRoom GetRoom
func GetRoom(ID int) (room *Room, errRet error) {
	Room := &Room{}
	o := orm.NewOrm()
	o.QueryTable("room").Filter("ID", ID).RelatedSel().One(Room)
	return Room, errRet
}

//CreateRoom _
func CreateRoom(Room Room) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&Room)
	o.Commit()
	return ID, err
}

//UpdateRoom _
func UpdateRoom(room Room) (errRet error) {
	o := orm.NewOrm()
	getUpdate, _ := GetRoom(room.ID)
	if getUpdate.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if getUpdate == nil {
		errRet = errors.New("ไม่พบข้อมูล")
	} else if errRet == nil {
		room.CreatedAt = getUpdate.CreatedAt
		room.Creator = getUpdate.Creator
		if num, errUpdate := o.Update(&room); errUpdate != nil {
			errRet = errUpdate
			_ = num
		}
	}
	return errRet
}

//DeleteRoom DeleteRoom
func DeleteRoom(ID int) (errRet error) {
	o := orm.NewOrm()
	unitDelete, _ := GetRoom(ID)
	if unitDelete.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if num, errDelete := o.Delete(&Room{ID: ID}); errDelete != nil && errRet == nil {
		errRet = errDelete
		_ = num
	}
	return errRet
}
