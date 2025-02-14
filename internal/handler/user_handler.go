package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tsaqiffatih/minddrift-server/config"
	"github.com/tsaqiffatih/minddrift-server/internal/dto"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
	"github.com/tsaqiffatih/minddrift-server/internal/service"
	"github.com/tsaqiffatih/minddrift-server/pkg/utils"
)

type UserHandler interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
	VerifyEmail(c *gin.Context)
	ResendEmail(c *gin.Context)
	ResetPassword(c *gin.Context)        //Change The Password
	ValidateResetToken(c *gin.Context)   //Validate Token Reset Password
	RequestResetPassword(c *gin.Context) //Request Reset Password and send email

	EnableTwoFA(c *gin.Context)
	DisableTwoFA(c *gin.Context)
	UpdateUserProfile(c *gin.Context)
	ChangeUserRole(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
	cfg         *config.Config
}

func NewUserHandler(userService service.UserService, cfg *config.Config) UserHandler {
	return &userHandler{
		cfg:         cfg,
		userService: userService,
	}
}

// **Register User**
func (h *userHandler) RegisterUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  validationErrors,
		})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  validationErrors,
		})
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     model.UserRole(req.Role),
	}

	createdUser, err := h.userService.RegisterUser(user)
	if err != nil {
		if err == model.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"errors":  "Email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  "Failed to register user",
		})
		return
	}

	response := dto.UserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Role:     string(createdUser.Role),
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully. Please check your email for verification.",
		"data":    gin.H{"user": response},
	})
}

// **Login User**
func (h *userHandler) LoginUser(c *gin.Context) {
	var loginData dto.LoginUserRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	if err := utils.ValidateStruct(loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	token, err := h.userService.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"errors":  err.Error(),
		})
		return
	}

	log.Println("Token:", token)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login success",
		"data":    gin.H{"token": token},
	})
}

// **Verification Email**
func (h *userHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Token is required",
		})
		return
	}

	err := h.userService.VerifyEmail(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email verified successfully",
	})
	log.Println("Email Verified Successfully")
}

func (h *userHandler) ResendEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Email is required",
		})
		return
	}

	err := h.userService.ResendEmail(req.Email)
	if err != nil {

		if err.Error() == "Invalid Email" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err,
			})
			return
		} else if err.Error() == "Email has been verified" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email verification sent successfully",
	})
}

// ResetPassword implements UserHandler.
func (h *userHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPassword

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  err.Error()})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	err := h.userService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Your password has been changed, please login with your new password",
	})
}

// ValidateResetToken implements UserHandler.
func (h *userHandler) ValidateResetToken(c *gin.Context) {
	var req dto.ValidateResetToken

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  err.Error()})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  validationErrors,
		})
		return
	}

	if err := h.userService.ValidateTokenResetPassword(req.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token valid",
	})
}

// **Reset Password**
func (h *userHandler) RequestResetPassword(c *gin.Context) {
	var resetData dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&resetData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  err.Error()})
		return
	}

	err := h.userService.RequestResetPassword(resetData.Email)
	if err != nil {
		if err.Error() == "Invalid Email" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"errors":  "Invalid Email, Your email is not registered",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"errors":  "Something went wrong, please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset link has been sent to your email",
	})
}

// **Activate 2FA**
func (h *userHandler) EnableTwoFA(c *gin.Context) {
	var request struct {
		Secret string `json:"secret"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "ID tidak valid"})
		return
	}

	if err := h.userService.EnableTwoFA(userID, request.Secret); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA berhasil diaktifkan"})
}

// **Unactivate 2FA**
func (h *userHandler) DisableTwoFA(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "ID tidak valid"})
		return
	}

	if err := h.userService.DisableTwoFA(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA berhasil dinonaktifkan"})
}

// **Update Profil**
func (h *userHandler) UpdateUserProfile(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateUserProfile(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// **Change User Role**
func (h *userHandler) ChangeUserRole(c *gin.Context) {
	var roleData struct {
		AdminID string         `json:"admin_id"`
		UserID  string         `json:"user_id"`
		NewRole model.UserRole `json:"new_role"`
	}

	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	adminID, err := uuid.Parse(roleData.AdminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "ID admin tidak valid"})
		return
	}

	userID, err := uuid.Parse(roleData.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "ID pengguna tidak valid"})
		return
	}

	if err := h.userService.ChangeUserRole(adminID, userID, roleData.NewRole); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peran pengguna berhasil diubah"})
}

// **Delete User**
func (h *userHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "ID tidak valid"})
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}
