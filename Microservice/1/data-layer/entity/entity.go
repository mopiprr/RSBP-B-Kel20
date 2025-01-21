package entity

import (
	"errors"
	"time"

	helpers "github.com/mci-its/backend-service/helpers/function"
	"gorm.io/gorm"
)

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type User struct {
	Pkid           int64         `gorm:"type:bigint;primary_key;autoIncrement" json:"pkid"`
	Username       string        `gorm:"not null;unique"                       json:"username"`
	Email          string        `gorm:"not null;unique"                       json:"email"`
	Password       string        `gorm:"not null"                              json:"password"`
	Nik            string        `gorm:"not null;unique;size:16;"              json:"nik"`
	FirstName      string        `gorm:"not null"                              json:"firstName"`
	LastName       string        `gorm:"not null"                              json:"lastName"`
	Gender         Gender        `gorm:"type:varchar(7);not null"              json:"gender"`
	Avatar         string        `                                             json:"avatar"`
	IsVerified     bool          `gorm:"default:false"                         json:"is_verified"`
	LoginAttempts  int           `gorm:"not null"                              json:"login_attempts"`
	SuspendedUntil time.Time     `gorm:"not null"                              json:"suspended_until"`
	Otp            string        `gorm:"null"                                  json:"otp"`
	OtpSent        int           `gorm:"not null;default:0"                    json:"otp_sent"`
	RolePkid       int64         `gorm:"not null"                              json:"role_pkid"`
	Role           Role          `gorm:"foreignKey:RolePkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"role"`
	OfficePkid     int64         `gorm:"not null"                              json:"office_pkid"`
	Office         Office        `gorm:"foreignKey:OfficePkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"office"`
	Logs           []LogActivity `gorm:"foreignKey:UserPkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"logs"`
	Contacts       []Contact     `gorm:"foreignKey:UserPkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"contacts"`
	CreatedBy      string        `                                             json:"created_by"`
	CreatedHost    string        `                                             json:"created_host"`
	UpdatedBy      string        `                                             json:"updated_by"`
	UpdatedHost    string        `                                             json:"updated_host"`
	DeletedBy      string        `                                             json:"deleted_by"`
	DeletedHost    string        `                                             json:"deleted_host"`
	IsDeleted      bool          `gorm:"default:false"                         json:"is_deleted"`

	Timestamp
}

type RoleName string

const (
	Superadmin  RoleName = "Superadmin"
	Supervisor  RoleName = "Supervisor"
	OfficeAdmin RoleName = "Office Admin"
)

type Role struct {
	Pkid        int64    `gorm:"primaryKey;autoIncrement"`
	RoleName    RoleName `gorm:"type:varchar(20);not null;unique"`
	Users       []User   `gorm:"foreignKey:RolePkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"users"`
	Offices     []Office `gorm:"many2many:role_in_offices;" json:"offices"`
	CreatedBy   string   `                                             json:"created_by"`
	CreatedHost string   `                                             json:"created_host"`
	UpdatedBy   string   `                                             json:"updated_by"`
	UpdatedHost string   `                                             json:"updated_host"`
	DeletedBy   string   `                                             json:"deleted_by"`
	DeletedHost string   `                                             json:"deleted_host"`
	IsDeleted   bool     `gorm:"default:false"                         json:"is_deleted"`

	Timestamp
}

type ActionType string

const (
	Create ActionType = "Create"
	Update ActionType = "Update"
	Delete ActionType = "Delete"
)

type LogActivity struct {
	Pkid        int64      `gorm:"type:bigint;primary_key;autoIncrement" json:"pkid"`
	UserPkid    int64      `gorm:"not null"                              json:"user_pkid"`
	User        User       `gorm:"foreignKey:UserPkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	IpAddress   string     `gorm:"not null"                              json:"ip_address"`
	ActionType  ActionType `gorm:"type:varchar(10);not null"             json:"activity_type"`
	ActionTime  time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"    json:"activity_time"`
	TargetTable string     `gorm:"not null"                              json:"target_table"`
	TargetPkid  int64      `gorm:"type:bigint;not null"                  json:"target_pkid"`
	OfficePkid  int64      `gorm:"not null"                              json:"office_pkid"`
	Office      Office     `gorm:"foreignKey:OfficePkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"office"`
	CreatedBy   string     `                                             json:"created_by"`
	CreatedHost string     `                                             json:"created_host"`
	UpdatedBy   string     `                                             json:"updated_by"`
	UpdatedHost string     `                                             json:"updated_host"`
	DeletedBy   string     `                                             json:"deleted_by"`
	DeletedHost string     `                                             json:"deleted_host"`
	IsDeleted   bool       `gorm:"default:false"                         json:"is_deleted"`

	Timestamp
}

type Contact struct {
	Pkid         int64  `gorm:"type:bigint;primary_key;autoIncrement" json:"pkid"`
	UserPkid     int64  `gorm:"not null"                              json:"user_pkid"`
	User         User   `gorm:"foreignKey:UserPkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	ContactValue string `gorm:"not null"                              json:"contact_value"`
	OfficePkid   int64  `gorm:"not null"                              json:"office_pkid"`
	Office       Office `gorm:"foreignKey:OfficePkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"office"`
	CreatedBy    string `                                             json:"created_by"`
	CreatedHost  string `                                             json:"created_host"`
	UpdatedBy    string `                                             json:"updated_by"`
	UpdatedHost  string `                                             json:"updated_host"`
	DeletedBy    string `                                             json:"deleted_by"`
	DeletedHost  string `                                             json:"deleted_host"`
	IsDeleted    bool   `gorm:"default:false"                         json:"is_deleted"`

	Timestamp
}

type Office struct {
	Pkid          int64  `gorm:"type:bigint;primary_key;autoIncrement" json:"pkid"`
	OfficeCode    string `gorm:"not null" json:"office_code"`
	OfficeName    string `gorm:"not null" json:"office_name"`
	OfficeAddress string `gorm:"not null" json:"office_address"`
	OfficePhone   string `gorm:"not null" json:"office_phone"`
	OfficeEmail   string `gorm:"not null" json:"office_email"`
	Roles         []Role `gorm:"many2many:role_in_offices;" json:"roles"`
	CreatedBy     string `                                             json:"created_by"`
	CreatedHost   string `                                             json:"created_host"`
	UpdatedBy     string `                                             json:"updated_by"`
	UpdatedHost   string `                                             json:"updated_host"`
	DeletedBy     string `                                             json:"deleted_by"`
	DeletedHost   string `                                             json:"deleted_host"`
	IsDeleted     bool   `gorm:"default:false"                         json:"is_deleted"`

	Timestamp
}

type RoleInOffice struct {
	RolePkid   int64  `gorm:"primaryKey"`
	OfficePkid int64  `gorm:"primaryKey"`
	Role       Role   `gorm:"foreignKey:RolePkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Office     Office `gorm:"foreignKey:OfficePkid;references:Pkid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}

	if len(u.Nik) != 16 {
		return errors.New("NIK must be exactly 16 characters long")
	}
	for _, char := range u.Nik {
		if char < '0' || char > '9' {
			return errors.New("NIK must contain only digits")
		}
	}

	return nil
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.RoleName != Superadmin && r.RoleName != Supervisor && r.RoleName != OfficeAdmin {
		return errors.New("invalid role name")
	}
	return nil
}
