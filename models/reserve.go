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
	Lecturer        string `orm:"null;size(300)"`
	Coordinate      string `orm:"null;size(300)"`
	MemberText      string `orm:"null;size(2000)"`
	MemberCount     int    `orm:"null"`
	DeviceAddOnText string `orm:"null;size(500)"`
	Status          int
	HideTitle       int
	HideFile        int
	DateBegin       time.Time         `form:"-" orm:"null;"`
	DateEnd         time.Time         `form:"-" orm:"null;"`
	Remark          string            `orm:"size(300)"`
	Room            *Room             `orm:"rel(fk)"`
	RoomName        string            `orm:"-"`
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

//RoomReserveResult *
type RoomReserveResult struct {
	ID              int
	Name            string
	SupportText     string
	LocationText    string
	AddOnDeviceText string
	RoomAdminText   string
	ImagePath1      string
	Remark          string
	Status          int
	ReserveNumber   int
	Title           string
	Coordinate      string
	HasReserve      bool
}

//RoomReserveList _
type RoomReserveList struct {
	ID              int
	DocNo           string
	Lock            bool
	Title           string
	Lecturer        string
	Coordinate      string
	MemberText      string
	MemberCount     int
	DeviceAddOnText string
	Status          int
	HideTitle       int
	HideFile        int
	DateBegin       string
	DateEnd         string
	Remark          string
	RoomName        string
	RoomID          int
}

//RoomReserveListJSON -
type RoomReserveListJSON struct {
	RoomReserveList *[]RoomReserveResult
	Paging          string
	Page            uint
	PerPage         uint
}

//MyReserveListJSON -
type MyReserveListJSON struct {
	MyReserveList *[]RoomReserveList
	Paging        string
	Page          uint
	PerPage       uint
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

//GetReserveFileList -
func GetReserveFileList(ID int) (rev *[]RoomReserveFile, errRet error) {
	RoomReserveFile := &[]RoomReserveFile{}
	o := orm.NewOrm()
	o.QueryTable("room_reserve_file").Filter("reserve_i_d", ID).RelatedSel().All(RoomReserveFile)
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

//GetReserveRoom -
func GetReserveRoom(ID int) (rev *RoomReserve, errRet error) {
	RoomReserve := &RoomReserve{}
	o := orm.NewOrm()
	o.QueryTable("room_reserve").Filter("ID", ID).RelatedSel().One(RoomReserve)
	if RoomReserve.ID != 0 {
		file, _ := GetReserveFileList(RoomReserve.ID)
		RoomReserve.RoomReserveFile = *file
	}
	return RoomReserve, errRet
}

//CreateReserveRoom _
func CreateReserveRoom(RoomReserve RoomReserve) (ID int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	ID, err = o.Insert(&RoomReserve)
	o.Commit()
	return ID, err
}

//UpdateReserveRoom _
func UpdateReserveRoom(RoomReserve RoomReserve) (err error) {
	room, _ := GetReserveRoom(RoomReserve.ID)
	RoomReserve.CreatedAt = room.CreatedAt
	RoomReserve.Creator = room.Creator
	o := orm.NewOrm()
	o.Begin()
	_, err = o.Update(&RoomReserve)
	o.Commit()
	return err
}

//GetReserveList -
func GetReserveList(currentPage, lineSize uint, term, status, room, beginDate, endDate string) (num int64, reserveListJSON []RoomReserveResult, err error) {
	o := orm.NewOrm()
	var sql = `SELECT
					room.i_d,
					room.name,
					room.support_text,
					room.location_text,
					room.add_on_device_text,
					room.room_admin_text,
					room.image_path1,
					room.remark,
					room.status,  					 
					(select reserve.i_d 
						from  room_reserve reserve 
						where reserve.room_id = room.i_d AND reserve.status = 1
						AND ('` + beginDate + `' BETWEEN reserve.date_begin AND reserve.date_end  
						OR '` + endDate + `'    BETWEEN reserve.date_begin AND reserve.date_end)    limit 1
					) AS reserve_number 
				FROM
					room  
				WHERE 1=1  `

	if status == "0" {
		sql += ` AND COALESCE((select reserve.i_d 
							from  room_reserve reserve 
							where reserve.room_id = room.i_d AND reserve.status = 1
							AND ('` + beginDate + `' BETWEEN reserve.date_begin AND reserve.date_end  
						    	OR '` + endDate + `'    BETWEEN reserve.date_begin AND reserve.date_end)    limit 1
						),0)  <> 1`
	}

	if room != "" {
		sql += ` AND room.i_d = ` + room + `  `
	}

	if status == "0" {
		sql += ` AND room.status = 1`
	}

	num, _ = o.Raw(sql).QueryRows(&reserveListJSON)
	if lineSize+currentPage > uint(num) {
		lineSize = uint(num)
	} else if currentPage > 0 {
		lineSize = lineSize + currentPage
	}
	if currentPage > lineSize {
		currentPage = 0
	}
	reserveListJSON = reserveListJSON[currentPage:lineSize]
	return num, reserveListJSON, err
}

//GetMyReserveList -
func GetMyReserveList(currentPage, lineSize uint, term, status, room, beginDate, endDate string) (num int64, myReserveListJSON []RoomReserveList, err error) {
	o := orm.NewOrm()
	var sql = ` select 
					   to_char(	T0.date_begin + interval '543' year, 'DD/MM/YYYY HH:MI') as date_begin, 
			           to_char(	T0.date_end + interval '543' year, 'DD/MM/YYYY HH:MI') as date_end,  
						*,
						T1.Name as room_name,
						T1.I_D as room_i_d 
				from  room_reserve T0 JOIN room T1 on T0.room_id = T1.i_d
				where 1=1  AND ('` + beginDate + `' BETWEEN T0.date_begin AND T0.date_end  
						   OR '` + endDate + `'    BETWEEN T0.date_begin AND T0.date_end)  `
	num, _ = o.Raw(sql).QueryRows(&myReserveListJSON)
	if lineSize+currentPage > uint(num) {
		lineSize = uint(num)
	} else if currentPage > 0 {
		lineSize = lineSize + currentPage
	}
	if currentPage > lineSize {
		currentPage = 0
	}
	myReserveListJSON = myReserveListJSON[currentPage:lineSize]
	return num, myReserveListJSON, err
}
