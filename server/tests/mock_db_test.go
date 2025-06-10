package templateapp

import (
	"errors"
	"gorm.io/gorm"
	"templateapp/models"
)

type MockDB struct {
	Users []models.User
}

func NewMockDBUser() *MockDB {
	return &MockDB{Users: []models.User{}}
}

func (db *MockDB) Create(value interface{}) *gorm.DB {
	user, ok := value.(*models.User)
	if !ok {
		return &gorm.DB{Error: errors.New("invalid type")}
	}

	for _, u := range db.Users {
		if u.Email == user.Email {
			return &gorm.DB{Error: errors.New("email already exists")}
		}
	}
	user.ID = int64(len(db.Users) + 1)
	db.Users = append(db.Users, *user)
	return &gorm.DB{Error: nil}
}

func (db *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	user, ok := dest.(*models.User)
	if !ok {
		return &gorm.DB{Error: errors.New("invalid type")}
	}

	email := conds[1].(string)
	for _, u := range db.Users {
		if u.Email == email {
			*user = u
			return &gorm.DB{Error: nil}
		}
	}
	return &gorm.DB{Error: gorm.ErrRecordNotFound}
}

func (db *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	// 추가 구현 가능
	return &gorm.DB{Error: nil}
}

func (db *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	// 추가 구현 가능
	return &gorm.DB{Error: nil}
}
