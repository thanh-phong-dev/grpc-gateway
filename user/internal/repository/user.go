package repository

import (
	"gorm.io/gorm"
	"user/internal/model"
)

type UserRepository interface {
	Create(user model.User) error
	FindOne(condition model.User) (*model.User, error)
}

type user struct {
	db *gorm.DB
}

func NewUserRepository(dbConn *gorm.DB) UserRepository {
	return &user{
		db: dbConn,
	}
}

func (u *user) Create(user model.User) error {
	return u.db.Create(&user).Error
}

func (u *user) FindOne(condition model.User) (*model.User, error) {
	user := model.User{}
	result := u.db.Where(condition).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
