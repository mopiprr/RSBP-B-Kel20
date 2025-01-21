package repository

import (
	"context"
	"math"

	"github.com/mci-its/backend-service/data-layer/entity"
	"github.com/mci-its/backend-service/helpers/dto"
	"gorm.io/gorm"
)

type (
	UserRepository interface {
		RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllUserRepositoryResponse, error)
		GetUserById(ctx context.Context, tx *gorm.DB, userPkid int64) (entity.User, error)
		GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error)
		CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, bool, error)
		CheckUsername(ctx context.Context, tx *gorm.DB, username string) (entity.User, bool, error)
		CheckNik(ctx context.Context, tx *gorm.DB, nik string) (entity.User, bool, error)
		UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		DeleteUser(ctx context.Context, tx *gorm.DB, userPkid int64) error
		GetUserForChangePassword(ctx context.Context, tx *gorm.DB, userPkid int64) (string, error)
		UpdateUserPassword(ctx context.Context, tx *gorm.DB, newPassword string) error
		UpdateLoginAttempt(ctx context.Context, tx *gorm.DB, email string, user entity.User) error
		SoftDeleteUser(ctx context.Context, tx *gorm.DB, userPkid int64) error
	}

	userRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllUserRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var users []entity.User
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if err := tx.WithContext(ctx).Model(&entity.User{}).Count(&count).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	if err := tx.WithContext(ctx).Scopes(Paginate(req.Page, req.PerPage)).Find(&users).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.GetAllUserRepositoryResponse{
		Users: users,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

func (r *userRepository) GetUserById(ctx context.Context, tx *gorm.DB, userPkid int64) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("Pkid = ?", userPkid).Take(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}

// ADD: CheckUsername
func (r *userRepository) CheckUsername(ctx context.Context, tx *gorm.DB, username string) (entity.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("username = ?", username).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}

// ADD: CheckNik
func (r *userRepository) CheckNik(ctx context.Context, tx *gorm.DB, nik string) (entity.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("nik = ?", nik).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, tx *gorm.DB, userPkid int64) error {
	if tx == nil {
		tx = r.db
	}

	// Hard delete: Permanently remove the user record from the database
	if err := tx.WithContext(ctx).Unscoped().Delete(&entity.User{}, "pkid = ?", userPkid).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserPassword(ctx context.Context, tx *gorm.DB, newPassword string) error {
	if tx == nil {
		tx = r.db
	}

	// Hanya update password
	if err := tx.WithContext(ctx).Update("password", newPassword).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserForChangePassword(ctx context.Context, tx *gorm.DB, userPkid int64) (string, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("pkid = ?", userPkid).Take(&user).Error; err != nil {
		return "", err
	}

	return user.Password, nil
}

func (r *userRepository) UpdateLoginAttempt(ctx context.Context, tx *gorm.DB, emailOrUsername string, user entity.User) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Model(&entity.User{}).Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).Updates(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) SoftDeleteUser(ctx context.Context, tx *gorm.DB, userPkid int64) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.User{}, "pkid = ?", userPkid).Error; err != nil {
		return err
	}

	return nil
}
