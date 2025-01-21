package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mci-its/backend-service/domain-layer/service"
	"github.com/mci-its/backend-service/helpers/dto"
	"github.com/mci-its/backend-service/helpers/utils"
)

type (
	UserController interface {
		Register(ctx *gin.Context)
		Login(ctx *gin.Context)
		Me(ctx *gin.Context)
		GetAllUser(ctx *gin.Context)
		SendVerificationEmail(ctx *gin.Context)
		VerifyEmail(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
		CheckToken(ctx *gin.Context)
		RefreshToken(ctx *gin.Context)
		GetUserByToken(ctx *gin.Context)
		ChangePassword(ctx *gin.Context)
		SendResetPassword(ctx *gin.Context)
		ResetPassword(ctx *gin.Context)
		VerifyOtp(ctx *gin.Context)
		ResendOtp(ctx *gin.Context)
		SoftDelete(ctx *gin.Context)
		UploadAvatar(ctx *gin.Context)
	}

	userController struct {
		userService service.UserService
	}
)

func NewUserController(us service.UserService) UserController {
	return &userController{
		userService: us,
	}
}

func (c *userController) Register(ctx *gin.Context) {
	var user dto.UserCreateRequest
	if err := ctx.ShouldBind(&user); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.RegisterUser(ctx.Request.Context(), user)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) GetAllUser(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.GetAllUserWithPagination(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	resp := utils.Response{
		Status:  true,
		Message: dto.MESSAGE_SUCCESS_GET_LIST_USER,
		Data:    result.Data,
		Meta:    result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *userController) Me(ctx *gin.Context) {
	userPkid := ctx.MustGet("user_pkid").(int64)

	result, err := c.userService.GetUserById(ctx.Request.Context(), userPkid)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	result, err := c.userService.Verify(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGIN, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) SendVerificationEmail(ctx *gin.Context) {
	var req dto.SendVerificationEmailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.userService.SendVerificationEmail(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) VerifyEmail(ctx *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.VerifyEmail(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VERIFY_EMAIL, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_VERIFY_EMAIL, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Update(ctx *gin.Context) {
	var req dto.UserUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userPkid := ctx.MustGet("user_pkid").(int64)
	result, err := c.userService.UpdateUser(ctx.Request.Context(), req, userPkid)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Delete(ctx *gin.Context) {
	userPkid := ctx.MustGet("user_pkid").(int64)

	if err := c.userService.DeleteUser(ctx.Request.Context(), userPkid); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, nil)
	ctx.JSON(http.StatusOK, res)
}
func (c *userController) CheckToken(ctx *gin.Context) {
	// Ambil token dari Context yang sudah diset oleh middleware
	token := ctx.GetString("token")
	if token == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, "Token not found", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Gunakan service untuk cek validitas token
	result, err := c.userService.CheckToken(ctx.Request.Context(), token)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess("Token Valid", result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) RefreshToken(ctx *gin.Context) {
	// Ambil token dari Context yang sudah diset oleh middleware
	token := ctx.GetString("token")
	if token == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, "Token not found", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Gunakan service untuk refresh token
	result, err := c.userService.RefreshToken(ctx.Request.Context(), token)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess("Token Refreshed", result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) GetUserByToken(ctx *gin.Context) {
	token := ctx.GetString("token")

	if token == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, "Token not found", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.GetUserByToken(ctx.Request.Context(), token)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess("Get User By Token", result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) ChangePassword(ctx *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userPkid := ctx.MustGet("user_pkid").(int64)
	err := c.userService.ChangePassword(ctx.Request.Context(), req, userPkid)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CHANGE_PASSWORD, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CHANGE_PASSWORD, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) SendResetPassword(ctx *gin.Context) {
	var req dto.SendResetPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.userService.SendResetPassword(ctx.Request.Context(), req.Email)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_SENT_RESET_PASSWORD, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SEND_RESET_PASSWORD, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.userService.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CHANGE_PASSWORD, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CHANGE_PASSWORD, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) VerifyOtp(ctx *gin.Context) {
	var req dto.VerifyOtpRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res, err := c.userService.VerifyOtp(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result := utils.BuildResponseSuccess("OTP Verified", res)
	ctx.JSON(http.StatusOK, result)
}

func (c *userController) ResendOtp(ctx *gin.Context) {
	var req dto.ResendOtpRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res, err := c.userService.ResendOtp(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result := utils.BuildResponseSuccess("OTP Resend", res)
	ctx.JSON(http.StatusOK, result)
}

func (c *userController) SoftDelete(ctx *gin.Context) {
	userPkid := ctx.MustGet("user_pkid").(int64)

	err := c.userService.SoftDeleteUser(ctx.Request.Context(), userPkid)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, nil)
	ctx.JSON(http.StatusOK, result)
}

func (c *userController) UploadAvatar(ctx *gin.Context) {
	userPkid := ctx.MustGet("user_pkid").(int64)

	file, err := ctx.FormFile("avatar")
	if err != nil {
		res := utils.BuildResponseFailed("Failed to get file", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res, err := c.userService.UploadAvatar(ctx.Request.Context(), file, userPkid)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPLOAD_AVATAR, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_AVATAR, res)
	ctx.JSON(http.StatusOK, result)
}
