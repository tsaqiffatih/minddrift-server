package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tsaqiffatih/minddrift-server/config"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
	"github.com/tsaqiffatih/minddrift-server/internal/repository"
	"github.com/tsaqiffatih/minddrift-server/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(user *model.User) (*model.User, error)
	LoginUser(email, password string) (string, error)
	VerifyEmail(token string) error
	ResendEmail(email string) error
	ResetPassword(email, newPassword string) error //Change The Password
	ValidateTokenResetPassword(token string) error //Validate Token Reset Password
	RequestResetPassword(email string) error       //Request Reset Password and send email
	EnableTwoFA(userID uuid.UUID, secret string) error
	DisableTwoFA(userID uuid.UUID) error
	UpdateUserProfile(user *model.User) (*model.User, error)
	ChangeUserRole(adminID, userID uuid.UUID, newRole model.UserRole) error
	DeleteUser(userID uuid.UUID) error
}

type userService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		repo: repo,
		cfg:  cfg,
	}
}

// **Register User**
func (s *userService) RegisterUser(user *model.User) (*model.User, error) {
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if existingUser != nil {
		return nil, model.ErrEmailAlreadyExists
	}

	err = utils.ValidatePasswordStrength(user.Password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	newUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	verificationToken, err := utils.GenerateTokenVerification(newUser.ID.String(), 24*time.Hour)
	if err != nil {
		return nil, err
	}

	go func() {
		verificationData := utils.EmailData{
			Username:         newUser.Username,
			VerificationLink: fmt.Sprintf("%s/verify-email?token=%s", s.cfg.FrontendURL, verificationToken),
			MindDriftEmail:   s.cfg.MindDriftEmail,
		}

		body, err := utils.GenerateEmailBody(utils.EmailVerification, verificationData)
		if err != nil {
			log.Println("Error generating email body:", err)
			return
		}

		err = utils.SendEmail(s.cfg, newUser.Email, "Email Verification", body)
		if err != nil {
			log.Printf("Error sending email for user %s: %v", newUser.Email, err)
		}
	}()

	return newUser, nil
}

// **Login User**
func (s *userService) LoginUser(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("Invalid email or password")
		}
		log.Println("Error getting user by email:", err)
		return "", errors.New("Something went wrong, please try again later")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("Invalid email or password")
	}

	if !user.EmailVerified {
		return "", errors.New("Email has not been verified. Please check your email.")
	}

	token, err := utils.GenerateJWT(s.cfg, user.ID, user.Role)
	if err != nil {
		log.Println("Error generating token:", err)
		return "", errors.New("Failed to generate token")
		// return "", errors.New("Internal server error")
	}

	return token, nil
}

// **Verification Email**
func (s *userService) VerifyEmail(token string) error {
	emailToken, err := utils.ParseTokenVerification(token)
	if err != nil {
		return err
	}

	log.Println("Email token:", emailToken)

	parsedEmailToken, err := uuid.Parse(emailToken)
	if err != nil {
		return errors.New("Invalid token format")
	}
	user, err := s.repo.GetUserByID(parsedEmailToken)
	if err != nil {
		log.Println("Error getting user by verification token:", err)
		return errors.New("Invalid token or token has expired")
	}

	if user == nil {
		log.Println("User not found by verification token")
		return errors.New("Invalid token or token has expired")
	}

	user.EmailVerified = true

	err = s.repo.UpdateUser(user)
	if err != nil {
		return errors.New("Failed to verify email")
	}

	return nil
}

// **Resend Email**
func (s *userService) ResendEmail(email string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if user == nil {
		return errors.New("Invalid Email")
	}

	if user.EmailVerified {
		return errors.New("Email has been verified")
	}

	verificationToken, err := utils.GenerateTokenVerification(user.ID.String(), 24*time.Hour)
	if err != nil {
		return err
	}

	go func() {
		verificationData := utils.EmailData{
			Username:         user.Username,
			VerificationLink: fmt.Sprintf("%s/verify-email?token=%s", s.cfg.FrontendURL, verificationToken),
			MindDriftEmail:   s.cfg.MindDriftEmail,
		}

		body, err := utils.GenerateEmailBody(utils.EmailVerification, verificationData)
		if err != nil {
			log.Println("Error generating email body:", err)
			return
		}

		err = utils.SendEmail(s.cfg, user.Email, "Email Verification", body)
		if err != nil {
			log.Println("Error sending email:", err)
		}
	}()

	return nil
}

// RequestResetPassword implements UserService.
func (s *userService) RequestResetPassword(email string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if user == nil {
		log.Println("Request reset password for non-existing email:", email)
		return nil
	}

	token, err := utils.GenerateTokenVerification(user.ID.String(), 30*time.Minute)
	if err != nil {
		return err
	}

	go func() {
		resetData := utils.EmailData{
			Username:         user.Username,
			VerificationLink: fmt.Sprintf("%s/verify-email?token=%s", s.cfg.FrontendURL, token),
			MindDriftEmail:   s.cfg.MindDriftEmail,
		}

		body, err := utils.GenerateEmailBody(utils.EmailResetPassword, resetData)
		if err != nil {
			log.Println("Error generating reset password email:", err)
			return
		}

		err = utils.SendEmail(s.cfg, user.Email, "Reset Password", body)
		if err != nil {
			log.Println("Error sending reset password email:", err)
		}
	}()

	return nil
}

// ValidateTokenResetPassword implements UserService.
func (s *userService) ValidateTokenResetPassword(token string) error {
	userID, err := utils.ParseTokenVerification(token)
	if err != nil {
		return err
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("Invalid Token, Your token is not registered")
	}

	_, err = s.repo.GetUserByID(parsedUserID)
	if err != nil {
		return errors.New("Invalid Token, Your token is not registered")
	}

	return nil
}

// **Reset Password**
func (s *userService) ResetPassword(token, newPassword string) error {
	userID, err := utils.ParseTokenVerification(token)
	if err != nil {
		return err
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("Invalid Token, Your token is not registered")
	}

	user, err := s.repo.GetUserByID(parsedUserID)
	if err != nil {
		return errors.New("Invalid Token, Your token is not registered")
	}

	if user == nil {
		return errors.New("Invalid Token, Your token is not registered")
	}

	err = utils.ValidatePasswordStrength(newPassword)
	if err != nil {
		return err
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
