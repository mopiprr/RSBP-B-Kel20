package repository

import (
	"context"

	"github.com/mci-its/backend-service/data-layer/entity"
	"gorm.io/gorm"
)

type (
	LogRepository interface {
		CreateLog(ctx context.Context, tx *gorm.DB, log entity.LogActivity) (entity.LogActivity, error)
		GetLog(ctx context.Context, tx *gorm.DB, userPkid int64) (entity.LogActivity, error)
	}

	logRepository struct {
		db *gorm.DB
	}
)

func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{
		db: db,
	}
}

func (r *logRepository) CreateLog(ctx context.Context, tx *gorm.DB, log entity.LogActivity) (entity.LogActivity, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&log).Error; err != nil {
		return entity.LogActivity{}, err
	}

	return log, nil
}

func (r *logRepository) GetLog(ctx context.Context, tx *gorm.DB, userPkid int64) (entity.LogActivity, error) {
	if tx == nil {
		tx = r.db
	}

	var log entity.LogActivity
	if err := tx.WithContext(ctx).Where("user_pkid = ?", userPkid).Order("created_at desc").Limit(10).Find(&log).Error; err != nil {
		return entity.LogActivity{}, err
	}

	return log, nil
}
