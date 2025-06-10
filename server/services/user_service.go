package services

import (
	"templateapp/database"
	"templateapp/models"
)

type UserService struct {
	DB database.Database
}

func (s *UserService) Register(email string, password string, name string) (models.User, error) {
	return models.User{}, nil
}

func (s *UserService) Login(email string, password string) (models.User, error) {
	return models.User{}, nil
}
