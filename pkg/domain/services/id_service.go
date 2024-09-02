package services

import "github.com/google/uuid"

type IDService interface {
	GenerateWorkflowID() string
}

type UUIDService struct{}

func NewUUIDService() IDService {
	return &UUIDService{}
}

func (s *UUIDService) GenerateWorkflowID() string {
	return uuid.New().String()
}
