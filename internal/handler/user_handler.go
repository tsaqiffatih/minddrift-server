package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tsaqiffatih/minddrift-server/internal/dto"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
	"github.com/tsaqiffatih/minddrift-server/internal/service"
	"github.com/tsaqiffatih/minddrift-server/pkg/utils"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// **Register User**
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.UserResponse{
		ID:            createdUser.ID,
		Username:      createdUser.Username,
		Email:         createdUser.Email,
		Role:          string(createdUser.Role),
		EmailVerified: createdUser.EmailVerified,
		TwoFAEnabled:  createdUser.TwoFAEnabled,
		CreatedAt:     createdUser.CreatedAt,
		UpdatedAt:     createdUser.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// **Login User**
func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// **Verification Email**
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.userService.VerifyEmail(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email berhasil diverifikasi"})
}

// **Reset Password**
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var resetData struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&resetData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.ResetPassword(resetData.Email, resetData.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil direset"})
}

// **Activate 2FA**
func (h *UserHandler) EnableTwoFA(c *gin.Context) {
	var request struct {
		Secret string `json:"secret"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.userService.EnableTwoFA(userID, request.Secret); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA berhasil diaktifkan"})
}

// **Unactivate 2FA**
func (h *UserHandler) DisableTwoFA(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.userService.DisableTwoFA(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA berhasil dinonaktifkan"})
}

// **Update Profil**
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateUserProfile(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// **Change User Role**
func (h *UserHandler) ChangeUserRole(c *gin.Context) {
	var roleData struct {
		AdminID string         `json:"admin_id"`
		UserID  string         `json:"user_id"`
		NewRole model.UserRole `json:"new_role"`
	}

	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, err := uuid.Parse(roleData.AdminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID admin tidak valid"})
		return
	}

	userID, err := uuid.Parse(roleData.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	if err := h.userService.ChangeUserRole(adminID, userID, roleData.NewRole); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peran pengguna berhasil diubah"})
}

// **Delete User**
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}
