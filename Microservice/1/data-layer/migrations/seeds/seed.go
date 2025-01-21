package seeds

import (
	"time"

	"github.com/mci-its/backend-service/data-layer/entity"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []entity.Role{
		{RoleName: entity.Superadmin, CreatedBy: "system", CreatedHost: "localhost"},
		{RoleName: entity.Supervisor, CreatedBy: "system", CreatedHost: "localhost"},
		{RoleName: entity.OfficeAdmin, CreatedBy: "system", CreatedHost: "localhost"},
	}

	for _, role := range roles {
		var existingRole entity.Role
		if err := db.Where("role_name = ?", role.RoleName).First(&existingRole).Error; err != nil {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func SeedOffices(db *gorm.DB) error {
	offices := []entity.Office{
		{
			OfficeCode:    "OFF001",
			OfficeName:    "Head Office",
			OfficeAddress: "123 Main St",
			OfficePhone:   "123-456-7890",
			OfficeEmail:   "headoffice@example.com",
			CreatedBy:     "system",
			CreatedHost:   "localhost",
		},
		{
			OfficeCode:    "OFF002",
			OfficeName:    "Branch Office",
			OfficeAddress: "456 Side St",
			OfficePhone:   "987-654-3210",
			OfficeEmail:   "branchoffice@example.com",
			CreatedBy:     "system",
			CreatedHost:   "localhost",
		},
	}

	for _, office := range offices {
		var existingOffice entity.Office
		if err := db.Where("office_code = ?", office.OfficeCode).First(&existingOffice).Error; err != nil {
			if err := db.Create(&office).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func SeedUsers(db *gorm.DB) error {
	users := []entity.User{
		{
			Username:    "superadmin",
			Email:       "superadmin@example.com",
			Password:    "password123",
			Nik:         "1234567890123456",
			FirstName:   "Super",
			LastName:    "Admin",
			Gender:      entity.Male,
			RolePkid:    1,
			OfficePkid:  1,
			CreatedBy:   "system",
			CreatedHost: "localhost",
		},
		{
			Username:    "officeadmin",
			Email:       "officeadmin@example.com",
			Password:    "password123",
			Nik:         "9876543210123456",
			FirstName:   "Office",
			LastName:    "Admin",
			Gender:      entity.Female,
			RolePkid:    3,
			OfficePkid:  2,
			CreatedBy:   "system",
			CreatedHost: "localhost",
		},
	}

	for _, user := range users {
		var existingUser entity.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func SeedLogActivities(db *gorm.DB) error {
	logActivities := []entity.LogActivity{
		{
			UserPkid:    1,
			IpAddress:   "192.168.1.1",
			ActionType:  entity.Create,
			ActionTime:  time.Now(),
			TargetTable: "roles",
			TargetPkid:  1,
			OfficePkid:  1,
			CreatedBy:   "system",
			CreatedHost: "localhost",
		},
		{
			UserPkid:    2,
			IpAddress:   "192.168.1.2",
			ActionType:  entity.Update,
			ActionTime:  time.Now(),
			TargetTable: "users",
			TargetPkid:  2,
			OfficePkid:  2,
			CreatedBy:   "system",
			CreatedHost: "localhost",
		},
	}

	for _, logActivity := range logActivities {
		if err := db.Create(&logActivity).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedContacts(db *gorm.DB) error {
	contacts := []entity.Contact{
		{
			UserPkid:     1,
			ContactValue: "098237498223",
			OfficePkid:   1,
			CreatedBy:    "system",
			CreatedHost:  "localhost",
		},
		{
			UserPkid:     2,
			ContactValue: "0827838732622",
			OfficePkid:   2,
			CreatedBy:    "system",
			CreatedHost:  "localhost",
		},
	}

	for _, contact := range contacts {
		if err := db.Create(&contact).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedAll(db *gorm.DB) error {
	if err := SeedRoles(db); err != nil {
		return err
	}

	if err := SeedOffices(db); err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedLogActivities(db); err != nil {
		return err
	}

	if err := SeedContacts(db); err != nil {
		return err
	}

	return nil
}
