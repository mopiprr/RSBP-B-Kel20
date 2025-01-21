package dto

import "errors"

const (
	MESSAGE_FAILED_GET_OFFICE = "Failed to get office"

	MESSAGE_SUCCESS_GET_OFFICE = "Success get office"
)

var (
	ErrOfficeNotFound = errors.New("office not found")
)
