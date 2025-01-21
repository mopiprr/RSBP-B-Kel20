package dto

import "errors"

const (
	MESSAGE_FAILED_GET_ROLE = "Failed to get role"

	MESSAGE_SUCCESS_GET_ROLE = "Success get role"
)

var (
	ErrRoleNotFound = errors.New("role not found")
	ErrGetRole      = errors.New("error get role")
)

type RoleResponse struct {
	Pkid     int64  `json:"pkid"`
	RoleName string `json:"role_name"`
}
