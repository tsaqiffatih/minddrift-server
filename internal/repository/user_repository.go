package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id uuid.UUID) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(&user).Error
}

func (r *userRepository) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&model.User{}, id).Error
}
