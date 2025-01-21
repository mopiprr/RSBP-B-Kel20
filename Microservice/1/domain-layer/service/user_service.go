package service

import (
	"bytes"
	"context"
	"html/template"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/mci-its/backend-service/data-layer/config"
	"github.com/mci-its/backend-service/data-layer/entity"
	"github.com/mci-its/backend-service/data-layer/repository"
	dto "github.com/mci-its/backend-service/helpers/dto"
	function "github.com/mci-its/backend-service/helpers/function"
	"github.com/mci-its/backend-service/helpers/utils"
)

type (
	UserService interface {
		RegisterUser(ctx context.Context, req dto.UserCreateRequest) (dto.RegisterUserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error)
		GetUserById(ctx context.Context, userPkid int64) (dto.UserResponse, error)
		GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error)
		SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error
		VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error)
		UpdateUser(ctx context.Context, req dto.UserUpdateRequest, userPkid int64) (dto.UserUpdateResponse, error)
		DeleteUser(ctx context.Context, userPkid int64) error
		Verify(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
		CheckToken(ctx context.Context, token string) (dto.CheckTokenResponse, error)
		RefreshToken(ctx context.Context, token string) (dto.RefreshTokenResponse, error)
		GetUserByToken(ctx context.Context, token string) (dto.GetUserByTokenResponse, error)
		ChangePassword(ctx context.Context, req dto.ChangePasswordRequest, userPkid int64) error
		SendResetPassword(ctx context.Context, receiverEmail string) error
		ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
		VerifyOtp(ctx context.Context, req dto.VerifyOtpRequest) (dto.VerifyOtpResponse, error)
		MakeVerificationOTP(ctx context.Context, receiverEmail string) (map[string]string, error)
		ResendOtp(ctx context.Context, req dto.ResendOtpRequest) (dto.ResendOtpResponse, error)
		SoftDeleteUser(ctx context.Context, userPkid int64) error
		UploadAvatar(ctx context.Context, avatar *multipart.FileHeader, userPkid int64) (dto.UploadAvatarResponse, error)
	}

	userService struct {
		userRepo   repository.UserRepository
		roleRepo   repository.RoleRepository
		officeRepo repository.OfficeRepository
		jwtService JWTService
	}
)

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, officeRepo repository.OfficeRepository, jwtService JWTService) UserService {
	return &userService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		officeRepo: officeRepo,
		jwtService: jwtService,
	}
}

const (
	LOCAL_URL            = "http://localhost:3000"
	VERIFY_EMAIL_ROUTE   = "api/user/verify-email"
	RESET_PASSWORD_ROUTE = "api/user/reset-password"
)

// TO DO: Edit RegisterUser to return token
func (s *userService) RegisterUser(ctx context.Context, req dto.UserCreateRequest) (dto.RegisterUserResponse, error) {
	_, flag, _ := s.userRepo.CheckEmail(ctx, nil, req.Email)
	if flag {
		return dto.RegisterUserResponse{}, dto.ErrEmailAlreadyExists
	}

	// ADD: Check if username already exists
	_, flag, _ = s.userRepo.CheckUsername(ctx, nil, req.Username)
	if flag {
		return dto.RegisterUserResponse{}, dto.ErrUsernameAlreadyExists
	}

	// ADD: Check if nik already exists
	_, flag, _ = s.userRepo.CheckNik(ctx, nil, req.Nik)
	if flag {
		return dto.RegisterUserResponse{}, dto.ErrNikAlreadyExists
	}

	user := entity.User{
		Nik:         req.Nik,
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		Email:       req.Email,
		RolePkid:    req.RolePkid,
		OfficePkid:  req.OfficePkid,
		Password:    req.Password,
		CreatedBy:   req.Username,
		CreatedHost: function.GetIpAdress(),
		UpdatedBy:   "",
		UpdatedHost: "",
		DeletedBy:   "",
		DeletedHost: "",
		IsVerified:  false,
	}

	userReg, err := s.userRepo.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.RegisterUserResponse{}, dto.ErrCreateUser
	}

	draftEmail, err := s.MakeVerificationOTP(ctx, user.Email)
	if err != nil {
		return dto.RegisterUserResponse{}, err
	}

	expired := time.Now().UTC().Add(time.Minute * 3).Format("2006-01-02 15:04:05")
	plainText := req.Email + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return dto.RegisterUserResponse{}, err
	}

	err = utils.SendMail(userReg.Email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return dto.RegisterUserResponse{}, err
	}

	return dto.RegisterUserResponse{
		Token: token,
	}, nil
}

func MakeVerificationEmail(receiverEmail string) (map[string]string, error) {
	expired := time.Now().UTC().Add(time.Minute * 3).Format("2006-01-02 15:04:05")
	plainText := receiverEmail + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return nil, err
	}

	verifyLink := LOCAL_URL + "/" + VERIFY_EMAIL_ROUTE + "?token=" + token

	readHtml, err := os.ReadFile("utils/email-template/base_mail.html")
	if err != nil {
		return nil, err
	}

	data := struct {
		Email  string
		Verify string
	}{
		Email:  receiverEmail,
		Verify: verifyLink,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return nil, err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return nil, err
	}

	draftEmail := map[string]string{
		"subject": "DroneMEQ - Verify Your Email",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}

func (s *userService) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	user, err := s.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.ErrEmailNotFound
	}

	draftEmail, err := MakeVerificationEmail(user.Email)
	if err != nil {
		return err
	}

	err = utils.SendMail(user.Email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	decryptedToken, err := utils.AESDecrypt(req.Token)

	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if !strings.Contains(decryptedToken, "_") {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	decryptedTokenSplit := strings.Split(decryptedToken, "_")
	email := decryptedTokenSplit[0]
	expired := decryptedTokenSplit[1]

	now := time.Now().UTC()
	expiredTime, err := time.Parse("2006-01-02 15:04:05", expired)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if now.Sub(expiredTime) > (time.Minute * 3) {
		return dto.VerifyEmailResponse{
			Email:      email,
			IsVerified: false,
		}, dto.ErrTokenExpired
	}

	user, err := s.userRepo.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
	}

	if user.IsVerified {
		return dto.VerifyEmailResponse{}, dto.ErrAccountAlreadyVerified
	}

	updatedUser, err := s.userRepo.UpdateUser(ctx, nil, entity.User{
		Pkid:       user.Pkid,
		IsVerified: true,
	})

	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUpdateUser
	}

	return dto.VerifyEmailResponse{
		Email:      email,
		IsVerified: updatedUser.IsVerified,
	}, nil
}

func (s *userService) GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := s.userRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var roleMap = make(map[int64]string)
	roles, err := s.roleRepo.GetAllRole(ctx, nil)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	for _, role := range roles {
		roleMap[role.Pkid] = string(role.RoleName)
	}

	var officeMap = make(map[int64]string)
	offices, err := s.officeRepo.GetAllOffice(ctx, nil)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	for _, office := range offices {
		officeMap[office.Pkid] = office.OfficeName
	}

	var datas []dto.UserResponse
	for _, user := range dataWithPaginate.Users {
		userRole := roleMap[user.RolePkid]
		userOffice := officeMap[user.OfficePkid]
		data := dto.UserResponse{
			Pkid:        user.Pkid,
			Role:        userRole,
			Office:      userOffice,
			Nik:         user.Nik,
			Username:    user.Username,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Gender:      user.Gender,
			Avatar:      user.Avatar,
			IsVerified:  user.IsVerified,
			CreatedBy:   user.CreatedBy,
			CreatedHost: user.CreatedHost,
			UpdatedBy:   user.UpdatedBy,
			UpdatedHost: user.UpdatedHost,
			DeletedBy:   user.DeletedBy,
			DeletedHost: user.DeletedHost,
		}

		datas = append(datas, data)
	}

	return dto.UserPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *userService) GetUserById(ctx context.Context, userPkid int64) (dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userPkid)

	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserById
	}
	userRole, err := s.roleRepo.GetRoleByID(ctx, nil, user.RolePkid)

	if err != nil {
		return dto.UserResponse{}, dto.ErrRoleNotFound
	}

	userOffice, err := s.officeRepo.GetOfficeByID(ctx, nil, user.OfficePkid)
	if err != nil {
		return dto.UserResponse{}, dto.ErrOfficeNotFound
	}

	return dto.UserResponse{
		Pkid:        user.Pkid,
		Nik:         user.Nik,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Gender:      user.Gender,
		Avatar:      user.Avatar,
		Role:        string(userRole.RoleName),
		Office:      userOffice.OfficeName,
		IsVerified:  user.IsVerified,
		CreatedBy:   user.CreatedBy,
		CreatedHost: user.CreatedHost,
		UpdatedBy:   user.UpdatedBy,
		UpdatedHost: user.UpdatedHost,
		DeletedBy:   user.DeletedBy,
		DeletedHost: user.DeletedHost,
	}, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserByEmail
	}

	return dto.UserResponse{
		Pkid:     user.Pkid,
		Nik:      user.Nik,
		Username: user.Username,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, req dto.UserUpdateRequest, userPkid int64) (dto.UserUpdateResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userPkid)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	data := entity.User{
		Pkid:        user.Pkid,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Avatar:      user.Avatar,
		Gender:      req.Gender,
		Email:       req.Email,
		Username:    req.Username,
		UpdatedBy:   req.Username,
		UpdatedHost: function.GetIpAdress(),
	}

	userUpdate, err := s.userRepo.UpdateUser(ctx, nil, data)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return dto.UserUpdateResponse{
		Pkid:      userUpdate.Pkid,
		FirstName: userUpdate.FirstName,
		LastName:  userUpdate.LastName,
		Gender:    userUpdate.Gender,
		Email:     userUpdate.Email,
		Username:  userUpdate.Username,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, userPkid int64) error {
	user, err := s.userRepo.GetUserById(ctx, nil, userPkid)
	if err != nil {
		return dto.ErrUserNotFound
	}

	err = s.userRepo.DeleteUser(ctx, nil, user.Pkid)
	if err != nil {
		return dto.ErrDeleteUser
	}

	return nil
}

func (s *userService) Verify(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	var check entity.User
	var flag bool
	var err error

	check, flag, err = s.userRepo.CheckEmail(ctx, nil, req.EmailOrUsername)
	if err != nil || !flag {
		check, flag, err = s.userRepo.CheckUsername(ctx, nil, req.EmailOrUsername)
		if err != nil || !flag {
			return dto.UserLoginResponse{}, dto.ErrEmailOrUsernameNotFound
		}
	}

	if time.Now().Before(check.SuspendedUntil) {
		return dto.UserLoginResponse{}, dto.ErrLoginSuspended
	}

	checkPassword, err := function.CheckPassword(check.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		check.LoginAttempts++

		if check.LoginAttempts >= 3 {
			check.SuspendedUntil = time.Now().Add(time.Minute * 10)
			check.LoginAttempts = 0
		}

		err = s.userRepo.UpdateLoginAttempt(ctx, nil, req.EmailOrUsername, entity.User{
			LoginAttempts:  check.LoginAttempts,
			SuspendedUntil: check.SuspendedUntil,
		})

		if err != nil {
			return dto.UserLoginResponse{}, err
		}

		return dto.UserLoginResponse{}, dto.ErrPasswordNotMatch
	}

	// ADD: If not verified force to verify
	if !check.IsVerified {
		return dto.UserLoginResponse{}, dto.ErrAccountNotVerified
	}

	token := s.jwtService.GenerateToken(strconv.FormatInt(check.Pkid, 10), strconv.FormatInt(check.RolePkid, 10), strconv.FormatInt(check.OfficePkid, 10))

	return dto.UserLoginResponse{
		Token: token,
	}, nil
}

func (s *userService) CheckToken(ctx context.Context, token string) (dto.CheckTokenResponse, error) {
	id, err := s.jwtService.GetUserPkidByToken(token)

	if err != nil {
		return dto.CheckTokenResponse{}, dto.ErrTokenInvalid
	}

	user, err := s.userRepo.GetUserById(ctx, nil, id)

	if err != nil {
		return dto.CheckTokenResponse{}, dto.ErrUserNotFound
	}

	return dto.CheckTokenResponse{
		UserPkid:   user.Pkid,
		Email:      user.Email,
		RolePkid:   user.RolePkid,
		OfficePkid: user.OfficePkid,
	}, nil
}

func (s *userService) RefreshToken(ctx context.Context, token string) (dto.RefreshTokenResponse, error) {
	userPkid, err := s.jwtService.GetUserPkidByToken(token)
	if err != nil {
		return dto.RefreshTokenResponse{}, dto.ErrTokenInvalid
	}

	user, err := s.userRepo.GetUserById(ctx, nil, userPkid)
	if err != nil {
		return dto.RefreshTokenResponse{}, dto.ErrUserNotFound
	}

	newToken := s.jwtService.GenerateTokenRefresh(
		strconv.FormatInt(user.Pkid, 10),
		strconv.FormatInt(user.RolePkid, 10),
		strconv.FormatInt(user.OfficePkid, 10),
	)

	return dto.RefreshTokenResponse{
		Token: newToken,
	}, nil
}

func (s *userService) GetUserByToken(ctx context.Context, token string) (dto.GetUserByTokenResponse, error) {
	userPkid, err := s.jwtService.GetUserPkidByToken(token)
	if err != nil {
		return dto.GetUserByTokenResponse{}, dto.ErrTokenInvalid
	}

	user, err := s.userRepo.GetUserById(ctx, nil, userPkid)
	if err != nil {
		return dto.GetUserByTokenResponse{}, dto.ErrUserNotFound
	}

	host := function.GetIpAdress()
	role, err := s.roleRepo.GetRoleByID(ctx, nil, user.RolePkid)

	if err != nil {
		return dto.GetUserByTokenResponse{}, dto.ErrRoleNotFound
	}

	return dto.GetUserByTokenResponse{
		UserPkid:    user.Pkid,
		Username:    user.Username,
		Email:       user.Email,
		Role:        string(role.RoleName),
		OfficePkid:  user.OfficePkid,
		Host:        host,
		CreatedBy:   user.CreatedBy,
		CreatedHost: user.CreatedHost,
		UpdatedBy:   user.UpdatedBy,
		UpdatedHost: user.UpdatedHost,
		DeletedBy:   user.DeletedBy,
		DeletedHost: user.DeletedHost,
		IsDeleted:   user.IsDeleted,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}, nil
}

func (s *userService) MakeVerificationOTP(ctx context.Context, receiverEmail string) (map[string]string, error) {
	otp := function.GenerateOTP()

	readHtml, err := os.ReadFile("helpers/utils/email-template/base_mail.html")
	if err != nil {
		return nil, err
	}

	data := struct {
		Email string
		Code  string
	}{
		Email: receiverEmail,
		Code:  otp,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return nil, err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByEmail(ctx, nil, receiverEmail)
	if err != nil {
		return nil, dto.ErrEmailNotFound
	}

	totalOtpSent := user.OtpSent + 1

	_, err = s.userRepo.UpdateUser(ctx, nil, entity.User{
		Pkid:    user.Pkid,
		Otp:     otp,
		OtpSent: totalOtpSent,
	})
	if err != nil {
		return nil, dto.ErrUpdateUser
	}

	draftEmail := map[string]string{
		"subject": "DroneMEQ - Verify Your Email",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}

// TO DO: Add function to verify OTP
func (s *userService) VerifyOtp(ctx context.Context, req dto.VerifyOtpRequest) (dto.VerifyOtpResponse, error) {
	decryptedToken, err := utils.AESDecrypt(req.Token)

	if err != nil {
		return dto.VerifyOtpResponse{}, dto.ErrTokenInvalid
	}

	if !strings.Contains(decryptedToken, "_") {
		return dto.VerifyOtpResponse{}, dto.ErrTokenInvalid
	}

	decryptedTokenSplit := strings.Split(decryptedToken, "_")
	email := decryptedTokenSplit[0]
	expired := decryptedTokenSplit[1]

	now := time.Now().UTC()
	expiredTime, err := time.Parse("2006-01-02 15:04:05", expired)
	if err != nil {
		return dto.VerifyOtpResponse{}, dto.ErrTokenInvalid
	}

	if now.Sub(expiredTime) > (time.Minute * 2) {
		return dto.VerifyOtpResponse{}, dto.ErrTokenExpired
	}

	user, err := s.userRepo.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.VerifyOtpResponse{}, dto.ErrEmailNotFound
	}

	if req.Otp != user.Otp {
		return dto.VerifyOtpResponse{}, dto.ErrOtpInvalid
	}

	_, err = s.userRepo.UpdateUser(ctx, nil, entity.User{
		Pkid:       user.Pkid,
		IsVerified: true,
		Otp:        "",
	})

	if err != nil {
		return dto.VerifyOtpResponse{}, dto.ErrUpdateUser
	}

	authToken := s.jwtService.GenerateToken(strconv.FormatInt(user.Pkid, 10), strconv.FormatInt(user.RolePkid, 10), strconv.FormatInt(user.OfficePkid, 10))

	return dto.VerifyOtpResponse{
		Token: authToken,
	}, nil
}

func (s *userService) ResendOtp(ctx context.Context, req dto.ResendOtpRequest) (dto.ResendOtpResponse, error) {
	decryptedToken, err := utils.AESDecrypt(req.Token)

	if err != nil {
		return dto.ResendOtpResponse{}, dto.ErrTokenInvalid
	}

	if !strings.Contains(decryptedToken, "_") {
		return dto.ResendOtpResponse{}, dto.ErrTokenInvalid
	}

	decryptedTokenSplit := strings.Split(decryptedToken, "_")
	email := decryptedTokenSplit[0]

	draftEmail, err := s.MakeVerificationOTP(ctx, email)
	if err != nil {
		return dto.ResendOtpResponse{}, err
	}

	err = utils.SendMail(email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return dto.ResendOtpResponse{}, err
	}

	expired := time.Now().UTC().Add(time.Minute * 3).Format("2006-01-02 15:04:05")
	plainText := email + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return dto.ResendOtpResponse{}, err
	}

	return dto.ResendOtpResponse{
		Token: token,
	}, nil

}

// TO DO: Add function to Change Password
func (s *userService) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest, userPkid int64) error {
	pwd, err := s.userRepo.GetUserForChangePassword(ctx, nil, userPkid)
	if err != nil {
		return err
	}

	checkPassword, err := function.CheckPassword(pwd, []byte(req.OldPassword))

	if err != nil || !checkPassword {
		return dto.ErrPasswordNotMatch
	}

	hashPassword, err := function.HashPassword(req.NewPassword)

	if err != nil {
		return err
	}

	_, err = s.userRepo.UpdateUser(ctx, nil, entity.User{
		Pkid:     userPkid,
		Password: hashPassword,
	})

	if err != nil {
		return err
	}

	return nil
}

// TO DO: Add function to Make Reset Password Link
func (s *userService) SendResetPassword(ctx context.Context, receiverEmail string) error {
	_, err := s.userRepo.GetUserByEmail(ctx, nil, receiverEmail)

	if err != nil {
		return dto.ErrEmailNotFound
	}

	expired := time.Now().UTC().Add(time.Minute * 2).Format("2006-01-02 15:04:05")
	plainText := receiverEmail + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return err
	}

	resetLink := LOCAL_URL + "/" + RESET_PASSWORD_ROUTE + "?token=" + token

	readHtml, err := os.ReadFile("helpers/utils/email-template/reset_password.html")
	if err != nil {
		return err
	}

	data := struct {
		Email string
		Link  string
	}{
		Email: receiverEmail,
		Link:  resetLink,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return err
	}

	draftEmail := map[string]string{
		"subject": "DroneMEQ - Reset Your Password",
		"body":    strMail.String(),
	}

	err = utils.SendMail(receiverEmail, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return err
	}

	return nil
}

// TO DO: Add function to Reset Password
func (s *userService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	decryptedToken, err := utils.AESDecrypt(req.Token)

	if err != nil {
		return dto.ErrTokenInvalid
	}

	if !strings.Contains(decryptedToken, "_") {
		return dto.ErrTokenInvalid
	}

	decryptedTokenSplit := strings.Split(decryptedToken, "_")
	email := decryptedTokenSplit[0]
	expired := decryptedTokenSplit[1]

	now := time.Now().UTC()
	expiredTime, err := time.Parse("2006-01-02 15:04:05", expired)
	if err != nil {
		return dto.ErrTokenInvalid
	}

	if now.Sub(expiredTime) > (time.Minute * 2) {
		return dto.ErrTokenExpired
	}

	user, err := s.userRepo.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.ErrUserNotFound
	}

	hashPassword, err := function.HashPassword(req.NewPassword)

	if err != nil {
		return err
	}

	_, err = s.userRepo.UpdateUser(ctx, nil, entity.User{
		Pkid:     user.Pkid,
		Password: hashPassword,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *userService) SoftDeleteUser(ctx context.Context, userPkid int64) error {
	user, err := s.userRepo.GetUserById(ctx, nil, userPkid)
	if err != nil {
		return dto.ErrUserNotFound
	}

	_, err = s.userRepo.UpdateUser(ctx, nil, entity.User{
		Pkid:        userPkid,
		IsDeleted:   true,
		DeletedBy:   user.Username,
		DeletedHost: function.GetIpAdress(),
	})

	if err != nil {
		return dto.ErrUpdateUser
	}

	err = s.userRepo.SoftDeleteUser(ctx, nil, user.Pkid)

	if err != nil {
		return dto.ErrDeleteUser
	}

	return nil
}

// ADD: function to Upload Avatar
func (s *userService) UploadAvatar(ctx context.Context, avatar *multipart.FileHeader, userPkid int64) (dto.UploadAvatarResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userPkid)
	if err != nil {
		return dto.UploadAvatarResponse{}, dto.ErrUserNotFound
	}

	fileHeader := avatar

	file, err := fileHeader.Open()

	if err != nil {
		return dto.UploadAvatarResponse{}, err
	}

	cldService, ctxBg := config.Credentials()

	res, err := cldService.Upload.Upload(ctxBg, file, uploader.UploadParams{
		PublicID:       "profile/" + user.Username,
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true)})

	if err != nil {
		return dto.UploadAvatarResponse{}, err
	}

	_, err = s.userRepo.UpdateUser(ctx, nil, entity.User{
		Pkid:        userPkid,
		Avatar:      res.SecureURL,
		UpdatedBy:   user.Username,
		UpdatedHost: function.GetIpAdress(),
	})

	if err != nil {
		return dto.UploadAvatarResponse{}, dto.ErrUpdateUser
	}

	return dto.UploadAvatarResponse{
		Avatar: res.SecureURL,
	}, nil
}
