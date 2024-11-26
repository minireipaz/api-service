package middlewares

import "minireipaz/pkg/domain/models"

func NewUnauthorizedError(message string) models.UnauthorizedError {
	return models.UnauthorizedError{Error: message}
}

func NewInvalidRequestError(message string, status int) models.InvalidRequestError {
	return models.InvalidRequestError{Error: message, Status: status}
}

func NewUnsupportedMediaTypeError(message string) models.UnsupportedMediaTypeError {
	return models.UnsupportedMediaTypeError{Error: message}
}

func NewTooManyRequestsError(message string) models.TooManyRequestsError {
	return models.TooManyRequestsError{Error: message}
}
