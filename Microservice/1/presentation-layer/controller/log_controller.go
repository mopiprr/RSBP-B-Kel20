package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mci-its/backend-service/domain-layer/service"
	"github.com/mci-its/backend-service/helpers/dto"
	"github.com/mci-its/backend-service/helpers/utils"
)

type (
	LogController interface {
		CreateLog(ctx *gin.Context)
		GetLogByUser(ctx *gin.Context)
	}

	logController struct {
		logService service.LogService
	}
)

func NewLogController(ls service.LogService) LogController {
	return &logController{
		logService: ls,
	}
}

func (c *logController) CreateLog(ctx *gin.Context) {
	var req dto.LogCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	userPkid := ctx.GetInt64("user_pkid")
	res, err := c.logService.CreateLog(ctx.Request.Context(), req, userPkid)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *logController) GetLogByUser(ctx *gin.Context) {
	userPkid := ctx.GetInt64("user_pkid")
	res, err := c.logService.GetLogByUser(ctx.Request.Context(), userPkid)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA, res)
	ctx.JSON(http.StatusOK, response)
}
