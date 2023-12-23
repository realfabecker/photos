package services

import (
	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
)

type UserService struct {
	UserRepository corpts.UserRepository
}

func NewUserService(r corpts.UserRepository) corpts.UserService {
	return &UserService{UserRepository: r}
}

func (s *UserService) GetUserByEmail(email string) (*cordom.User, error) {
	return s.UserRepository.GetUserByEmail(email)
}
