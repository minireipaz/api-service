package services

import "github.com/google/uuid"

type IDService interface {
    GenerateWorkflowID() uuid.UUID
}

type UUIDService struct{}

func NewUUIDService() IDService {
    return &UUIDService{}
}

func (s *UUIDService) GenerateWorkflowID() uuid.UUID {
    return uuid.New()
}
