package models

import (
	"time"
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
	DocNo       string `orm:"size(20)"`
	Lock        bool
	Status      int
	FilePath1   string    `orm:"size(300)"`
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
