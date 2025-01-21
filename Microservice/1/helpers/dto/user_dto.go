package dto

import (
	"errors"
	"mime/multipart"
	"time"

	"gorm.io/gorm"

	"github.com/mci-its/backend-service/data-layer/entity"
)

const (
	// Failed
	MESSAGE_FAILED_GET_DATA_FROM_BODY      = "failed get data from body"
	MESSAGE_FAILED_REGISTER_USER           = "failed create user"
	MESSAGE_FAILED_GET_LIST_USER           = "failed get list user"
	MESSAGE_FAILED_GET_USER_TOKEN          = "failed get user token"
	MESSAGE_FAILED_TOKEN_NOT_VALID         = "token not valid"
	MESSAGE_FAILED_TOKEN_NOT_FOUND         = "token not found"
	MESSAGE_FAILED_GET_USER                = "failed get user"
	MESSAGE_FAILED_LOGIN                   = "failed login"
	MESSAGE_FAILED_WRONG_EMAIL_OR_PASSWORD = "wrong email or password"
	MESSAGE_FAILED_UPDATE_USER             = "failed update user"
	MESSAGE_FAILED_DELETE_USER             = "failed delete user"
	MESSAGE_FAILED_PROSES_REQUEST          = "failed proses request"
	MESSAGE_FAILED_DENIED_ACCESS           = "denied access"
	MESSAGE_FAILED_VERIFY_EMAIL            = "failed verify email"
	MESSAGE_FAILED_CHANGE_PASSWORD         = "failed change password"
	MESSAGE_FAILED_SENT_RESET_PASSWORD     = "failed sent reset password link"
	MESSAGE_FAILED_UPLOAD_AVATAR           = "failed upload avatar"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER           = "success create user"
	MESSAGE_SUCCESS_GET_LIST_USER           = "success get list user"
	MESSAGE_SUCCESS_GET_USER                = "success get user"
	MESSAGE_SUCCESS_LOGIN                   = "success login"
	MESSAGE_SUCCESS_UPDATE_USER             = "success update user"
	MESSAGE_SUCCESS_DELETE_USER             = "success delete user"
	MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS = "success send verification email"
	MESSAGE_SUCCESS_VERIFY_EMAIL            = "success verify email"
	MESSAGE_SUCCESS_CHANGE_PASSWORD         = "success change password"
	MESSAGE_SEND_RESET_PASSWORD             = "success sent reset password link"
	MESSAGE_SUCCESS_GET_DATA                = "success get data"
	MESSAGE_SUCCESS_UPLOAD_AVATAR           = "success upload avatar"
)

var (
	ErrCreateUser              = errors.New("failed to create user")
	ErrGetAllUser              = errors.New("failed to get all user")
	ErrGetUserById             = errors.New("failed to get user by id")
	ErrGetUserByEmail          = errors.New("failed to get user by email")
	ErrEmailAlreadyExists      = errors.New("email already exist")
	ErrUsernameAlreadyExists   = errors.New("username already exist")
	ErrNikAlreadyExists        = errors.New("nik already exist")
	ErrUpdateUser              = errors.New("failed to update user")
	ErrUserNotAdmin            = errors.New("user not admin")
	ErrUserNotFound            = errors.New("user not found")
	ErrEmailNotFound           = errors.New("email not found")
	ErrEmailOrUsernameNotFound = errors.New("username or email not found")
	ErrDeleteUser              = errors.New("failed to delete user")
	ErrPasswordNotMatch        = errors.New("password not match")
	ErrEmailOrPassword         = errors.New("wrong email or password")
	ErrAccountNotVerified      = errors.New("account not verified")
	ErrTokenInvalid            = errors.New("token invalid")
	ErrTokenExpired            = errors.New("token expired")
	ErrAccountAlreadyVerified  = errors.New("account already verified")
	ErrLoginSuspended          = errors.New("login suspended")
	ErrOtpInvalid              = errors.New("otp invalid")
)

type (
	UserCreateRequest struct {
		Nik        string        `json:"nik"        form:"nik"        binding:"required,numeric,len=16"`
		Username   string        `json:"username"   form:"username"   binding:"required"`
		Email      string        `json:"email"      form:"email"      binding:"required,email"`
		Password   string        `json:"password"   form:"password"   binding:"required,min=8,containsany=0123456789,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*()"`
		ConfirmPwd string        `json:"confirm_pwd" form:"confirm_pwd" binding:"required,eqfield=Password"`
		FirstName  string        `json:"first_name" form:"first_name" binding:"required"`
		LastName   string        `json:"last_name"  form:"last_name"  binding:"required"`
		Gender     entity.Gender `json:"gender"     form:"gender"     binding:"required"`
		RolePkid   int64         `json:"role"       form:"role_pkid"       binding:"required"`
		OfficePkid int64         `json:"office" form:"office_pkid" binding:"required"`
	}

	UserResponse struct {
		Pkid        int64         `json:"pkid"`
		Nik         string        `json:"nik"`
		Username    string        `json:"username"`
		Email       string        `json:"email"`
		FirstName   string        `json:"first_name"`
		LastName    string        `json:"last_name"`
		Avatar      string        `json:"avatar"`
		Gender      entity.Gender `json:"gender"`
		Role        string        `json:"role"`
		Office      string        `json:"office"`
		CreatedBy   string        `json:"created_by" `
		CreatedHost string        `json:"created_host" `
		UpdatedBy   string        `json:"updated_by" `
		UpdatedHost string        `json:"updated_host" `
		DeletedBy   string        `json:"deleted_by"`
		DeletedHost string        `json:"deleted_host"`
		IsVerified  bool          `json:"is_verified"`
	}

	RegisterUserResponse struct {
		Token string `json:"token"`
	}

	UserPaginationResponse struct {
		Data []UserResponse `json:"data"`
		PaginationResponse
	}

	GetAllUserRepositoryResponse struct {
		Users []entity.User
		PaginationResponse
	}

	UserUpdateRequest struct {
		FirstName   string                `json:"first_name" form:"first_name" binding:"required"`
		LastName    string                `json:"last_name"  form:"last_name"  binding:"required"`
		Avatar      *multipart.FileHeader `json:"avatar,omitempty"     form:"avatar,omitempty"` // Not Required
		Gender      entity.Gender         `json:"gender"     form:"gender"     binding:"required"`
		Email       string                `json:"email"      form:"email"      binding:"required,email"`
		Username    string                `json:"username"   form:"username"   binding:"required"`
		UpdatedBy   string                `json:"updated_by"`
		UpdatedHost string                `json:"updated_host"`
	}

	UserUpdateResponse struct {
		Pkid      int64         `json:"pkid"`
		FirstName string        `json:"first_name"`
		LastName  string        `json:"last_name"`
		Avatar    string        `json:"avatar,omitempty"`
		Gender    entity.Gender `json:"gender"`
		Email     string        `json:"email"`
		Username  string        `json:"username"`
	}

	SendVerificationEmailRequest struct {
		Email string `json:"email" form:"email" binding:"required,email"`
	}

	VerifyEmailRequest struct {
		Token string `json:"token" form:"token" binding:"required"`
	}

	VerifyEmailResponse struct {
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
	}

	SendCodeVerificationRequest struct {
		Email string `json:"email" form:"email" binding:"required,email"`
	}

	UserLoginRequest struct {
		EmailOrUsername string `json:"email_or_username" form:"email_or_username" binding:"required"`
		Password        string `json:"password" form:"password" binding:"required"`
	}

	UserLoginResponse struct {
		Token string `json:"token"`
	}

	UpdateStatusIsVerifiedRequest struct {
		Pkid       int64 `json:"pkid" form:"pkid" binding:"required"`
		IsVerified bool  `json:"is_verified" form:"is_verified"`
	}

	CheckTokenResponse struct {
		UserPkid   int64  `json:"user_pkid"`
		Email      string `json:"email"`
		RolePkid   int64  `json:"role_pkid"`
		OfficePkid int64  `json:"office_pkid"`
	}

	RefreshTokenResponse struct {
		Token string `json:"token"`
	}

	GetUserByTokenResponse struct {
		UserPkid    int64          `json:"user_pkid"`
		Username    string         `json:"username"`
		Email       string         `json:"email"`
		Role        string         `json:"role"`
		OfficePkid  int64          `json:"office_pkid"`
		Host        string         `json:"host"`
		CreatedBy   string         `json:"created_by"`
		CreatedHost string         `json:"created_host"`
		UpdatedBy   string         `json:"updated_by"`
		UpdatedHost string         `json:"updated_host"`
		DeletedBy   string         `json:"deleted_by"`
		DeletedHost string         `json:"deleted_host"`
		IsDeleted   bool           `json:"is_deleted"`
		CreatedAt   time.Time      `json:"created_at"`
		UpdatedAt   time.Time      `json:"updated_at"`
		DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	}

	ChangePasswordRequest struct {
		OldPassword string `json:"old_password" form:"old_password" binding:"required"`
		NewPassword string `json:"new_password" form:"new_password" binding:"required,min=8,containsany=0123456789,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*()"`
		ConfirmPwd  string `json:"confirm_pwd" form:"confirm_pwd" binding:"required,eqfield=NewPassword"`
	}

	SendResetPasswordRequest struct {
		Email string `json:"email" form:"email" binding:"required,email"`
	}

	ResetPasswordRequest struct {
		Token       string `json:"token" form:"token" binding:"required"`
		NewPassword string `json:"new_password" form:"new_password" binding:"required,min=8,containsany=0123456789,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*()"`
		ConfirmPwd  string `json:"confirm_pwd" form:"confirm_pwd" binding:"required,eqfield=NewPassword"`
	}

	VerifyOtpRequest struct {
		Otp   string `json:"otp" form:"otp" binding:"required"`
		Token string `json:"token" form:"token" binding:"required"`
	}

	VerifyOtpResponse struct {
		Token string `json:"token"`
	}

	ResendOtpRequest struct {
		Token string `json:"token" form:"token" binding:"required"`
	}

	ResendOtpResponse struct {
		Token string `json:"token"`
	}

	UploadAvatarResponse struct {
		Avatar string `json:"avatar"`
	}
)
