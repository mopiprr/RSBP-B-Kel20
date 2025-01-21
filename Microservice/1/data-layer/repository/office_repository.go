package repository

import (
	"context"

	"github.com/mci-its/backend-service/data-layer/entity"
	"gorm.io/gorm"
)

type (
	OfficeRepository interface {
		GetAllOffice(ctx context.Context, tx *gorm.DB) ([]entity.Office, error)
		GetOfficeByID(ctx context.Context, tx *gorm.DB, officePkid int64) (entity.Office, error)
	}

	officeRepository struct {
		db *gorm.DB
	}
)

func NewOfficeRepository(db *gorm.DB) OfficeRepository {
	return &officeRepository{
		db: db,
	}
}

func (r *officeRepository) GetAllOffice(ctx context.Context, tx *gorm.DB) ([]entity.Office, error) {
	if tx == nil {
		tx = r.db
	}

	var offices []entity.Office
	if err := tx.WithContext(ctx).Find(&offices).Error; err != nil {
		return nil, err
	}

	return offices, nil
}

func (r *officeRepository) GetOfficeByID(ctx context.Context, tx *gorm.DB, officePkid int64) (entity.Office, error) {
	if tx == nil {
		tx = r.db
	}

	var office entity.Office
	if err := tx.WithContext(ctx).Where("pkid = ?", officePkid).First(&office).Error; err != nil {
		return entity.Office{}, err
	}

	return office, nil
}
