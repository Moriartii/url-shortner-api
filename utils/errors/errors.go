package errors

import (
	"errors"
	"net/http"
)

type RestErr struct {
	Message     string `json:"mesage"`
	Status      int    `json:"code"`
	Description string `json:"error"`
}

func NewError(msg string) error {
	return errors.New(msg)
}

func (rest *RestErr) Is(target *RestErr) bool {
	return rest.Message == target.Message
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusBadRequest,
		Description: "bad_request",
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusNotFound,
		Description: "not_found",
	}
}

func NewAlreadyExistError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusConflict,
		Description: "already_exist_error",
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusInternalServerError,
		Description: "internal_server_error",
	}
}

func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusUnauthorized,
		Description: "unauthorized",
	}
}
