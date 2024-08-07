package middlewares

import "minireipaz/pkg/domain/models"

func NewUnauthorizedError(message string) models.UnauthorizedError {
	return models.UnauthorizedError{Error: message}
}

func NewInvalidRequestError(message string) models.InvalidRequestError {
	return models.InvalidRequestError{Error: message}
}

func NewUnsupportedMediaTypeError(message string) models.UnsupportedMediaTypeError {
	return models.UnsupportedMediaTypeError{Error: message}
}

func NewTooManyRequestsError(message string) models.TooManyRequestsError {
	return models.TooManyRequestsError{Error: message}
}
