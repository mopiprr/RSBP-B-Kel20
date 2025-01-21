package repository

import (
	"context"

	"github.com/mci-its/backend-service/data-layer/entity"
	"gorm.io/gorm"
)

type (
	RoleRepository interface {
		GetRoleByID(ctx context.Context, tx *gorm.DB, rolePkid int64) (entity.Role, error)
		GetAllRole(ctx context.Context, tx *gorm.DB) ([]entity.Role, error)
	}

	roleRepository struct {
		db *gorm.DB
	}
)

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) GetRoleByID(ctx context.Context, tx *gorm.DB, rolePkid int64) (entity.Role, error) {
	if tx == nil {
		tx = r.db
	}

	var role entity.Role
	if err := tx.WithContext(ctx).Where("pkid = ?", rolePkid).First(&role).Error; err != nil {
		return entity.Role{}, err
	}

	return role, nil
}

func (r *roleRepository) GetAllRole(ctx context.Context, tx *gorm.DB) ([]entity.Role, error) {
	if tx == nil {
		tx = r.db
	}

	var roles []entity.Role
	if err := tx.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}
