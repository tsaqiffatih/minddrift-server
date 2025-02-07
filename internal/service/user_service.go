package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
	"github.com/tsaqiffatih/minddrift-server/internal/repository"
	"github.com/tsaqiffatih/minddrift-server/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(user *model.User) (*model.User, error)
	LoginUser(email, password string) (*model.User, error)
	VerifyEmail(userID uuid.UUID) error
	ResetPassword(email, newPassword string) error
	EnableTwoFA(userID uuid.UUID, secret string) error
	DisableTwoFA(userID uuid.UUID) error
	UpdateUserProfile(user *model.User) (*model.User, error)
	ChangeUserRole(adminID, userID uuid.UUID, newRole model.UserRole) error
	DeleteUser(userID uuid.UUID) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// **Register User**
func (s *userService) RegisterUser(user *model.User) (*model.User, error) {
	existingUser, _ := s.repo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	user.ID = uuid.New()

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// **Login User**
func (s *userService) LoginUser(email, password string) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	return user, nil
}

// **Verification Email**
func (s *userService) VerifyEmail(userID uuid.UUID) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	user.TwoFAEnabled = true
	return s.repo.UpdateUser(user)
}

// **Reset Password**
func (s *userService) ResetPassword(email, newPassword string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("email tidak ditemukan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.repo.UpdateUser(user)
}

// **Activated 2FA**
func (s *userService) EnableTwoFA(userID uuid.UUID, secret string) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	user.TwoFAEnabled = true
	user.TwoFASecret = secret
	return s.repo.UpdateUser(user)
}

// **Unactivated 2FA**
func (s *userService) DisableTwoFA(userID uuid.UUID) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	user.TwoFAEnabled = false
	user.TwoFASecret = ""
	return s.repo.UpdateUser(user)
}

// **Update User Profil**
func (s *userService) UpdateUserProfile(user *model.User) (*model.User, error) {
	existingUser, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.UpdatedAt = time.Now()

	err = s.repo.UpdateUser(existingUser)
	if err != nil {
		return nil, err
	}

	return existingUser, nil
}

// **Change Role User (Just Admin)**
func (s *userService) ChangeUserRole(adminID, userID uuid.UUID, newRole model.UserRole) error {
	admin, err := s.repo.GetUserByID(adminID)
	if err != nil {
		return errors.New("admin tidak ditemukan")
	}
	if admin.Role != model.Admin {
		return errors.New("hanya admin yang dapat mengubah peran pengguna")
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	user.Role = newRole
	return s.repo.UpdateUser(user)
}

// **Delete User**
func (s *userService) DeleteUser(userID uuid.UUID) error {
	return s.repo.DeleteUser(userID)
}
