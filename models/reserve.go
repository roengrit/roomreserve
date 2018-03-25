package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

//RoomReserve _
type RoomReserve struct {
	ID              int
	DocNo           string `orm:"size(20)"`
	Lock            bool
	Title           string `orm:"size(300)"`
	Lecturer        string `orm:"size(300)"`
	Coordinate      string `orm:"size(300)"`
	MemberText      string `orm:"size(2000)"`
	MemberCount     int
	DeviceAddOnText string `orm:"size(500)"`
	Status          int
	HideTitle       int
	HideFile        int
	DateBegin       time.Time
	DateEnd         time.Time
	Remark          string            `orm:"size(300)"`
	Creator         *User             `orm:"rel(fk)"`
	CreatedAt       time.Time         `orm:"auto_now_add;type(datetime)"`
	Editor          *User             `orm:"null;rel(fk)"`
	EditedAt        time.Time         `orm:"null;auto_now;type(datetime)"`
	RoomReserveFile []RoomReserveFile `orm:"-"`
}

//RoomReserveFile _
type RoomReserveFile struct {
	ID          int
	ReserveID   int
	Lock        bool
	Status      int
	FilePath1   string    `orm:"size(300)"`
	FileName    string    `orm:"size(300)"`
	DeleteFile1 int       `orm:"-"`
	Creator     *User     `orm:"rel(fk)"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	Editor      *User     `orm:"null;rel(fk)"`
	EditedAt    time.Time `orm:"null;auto_now;type(datetime)"`
}

//RoomReserveListJSON RoomReserveListJSON
type RoomReserveListJSON struct {
	RoomList *[]Room
	Paging   string
	Page     uint
	PerPage  uint
}

func init() {
	orm.RegisterModel(new(RoomReserve), new(RoomReserveFile))
}

//GetReserveFile -
func GetReserveFile(ID int) (rev *RoomReserveFile, errRet error) {
	RoomReserveFile := &RoomReserveFile{}
	o := orm.NewOrm()
	o.QueryTable("room_reserve_file").Filter("ID", ID).RelatedSel().One(RoomReserveFile)
	return RoomReserveFile, errRet
}

//CreateReserveFile _
func CreateReserveFile(RoomReserveFile RoomReserveFile) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&RoomReserveFile)
	o.Commit()
	return ID, err
}

//DeleteReserveFile -
func DeleteReserveFile(ID int) (errRet error) {
	o := orm.NewOrm()
	unitDelete, _ := GetRoom(ID)
	if unitDelete.Lock {
		errRet = errors.New("ข้อมูลถูก Lock ไม่สามารถแก้ไขได้")
	}
	if num, errDelete := o.Delete(&RoomReserveFile{ID: ID}); errDelete != nil && errRet == nil {
		errRet = errDelete
		_ = num
	}
	return errRet
}
