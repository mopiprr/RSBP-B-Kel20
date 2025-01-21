package dto

import (
	"errors"
	"time"

	"github.com/mci-its/backend-service/data-layer/entity"
)

const (
	// Failed
	MESSAGE_FAILED_CREATE_LOG = "failed to create log"
	MESSAGE_FAILED_DELETE_LOG = "failed to delete log"
	// Success
	MESSAGE_SUCCESS_CREATE_LOG = "success create log"
	MESSAGE_SUCCESS_DELETE_LOG = "success delete log"
)

var (
	ErrCreateLog = errors.New("failed to create log")
	ErrDeleteLog = errors.New("failed to delete log")
)

type (
	LogCreateRequest struct {
		ActionType  entity.ActionType `json:"activity_type" binding:"required"`
		ActionTime  time.Time         `json:"activity_time" binding:"required"`
		TargetTable string            `json:"target_table" binding:"required"`
		TargetPkid  int64             `json:"target_pkid" binding:"required"`
		CreatedBy   string            `json:"created_by"`
		UpdatedBy   string            `json:"updated_by"`
		UpdatedHost string            `json:"updated_host"`
		DeletedBy   string            `json:"deleted_by"`
		DeletedHost string            `json:"deleted_host"`
		IsDeleted   bool              `json:"is_deleted"`
	}

	LogCreateResponse struct {
		Pkid        int64             `json:"pkid"`
		UserPkid    int64             `json:"user_pkid"`
		IpAddress   string            `json:"ip_address"`
		OfficePkid  int64             `json:"office_pkid"`
		ActionType  entity.ActionType `json:"activity_type"`
		ActionTime  time.Time         `json:"activity_time"`
		TargetTable string            `json:"target_table"`
		TargetPkid  int64             `json:"target_pkid"`
	}

	GetLogResponse struct {
		UserPkid    int64             `json:"user_pkid"`
		IpAddress   string            `json:"ip_address"`
		OfficePkid  int64             `json:"office_pkid"`
		ActionType  entity.ActionType `json:"activity_type"`
		ActionTime  time.Time         `json:"activity_time"`
		TargetTable string            `json:"target_table"`
		TargetPkid  int64             `json:"target_pkid"`
	}
)
