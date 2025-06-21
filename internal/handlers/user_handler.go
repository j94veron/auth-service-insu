package handlers

import (
	"net/http"
	"strconv"

	"auth-aca/internal/models"
	"auth-aca/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo user.Repository
}

func NewUserHandler(userRepo user.Repository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

type CreateUserRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	Name           string `json:"name" binding:"required"`
	LastName       string `json:"lastName" binding:"required"`
	CommercialZone string `json:"commercialZone"`
	Warehouse      string `json:"warehouse"`
	RoleID         uint   `json:"roleId" binding:"required"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user := models.User{
		Email:          req.Email,
		Password:       string(hashedPassword),
		Name:           req.Name,
		LastName:       req.LastName,
		CommercialZone: req.CommercialZone,
		Warehouse:      req.Warehouse,
		RoleID:         req.RoleID,
	}

	if err := h.userRepo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

type UpdateUserRequest struct {
	Name           string `json:"name"`
	LastName       string `json:"lastName"`
	CommercialZone string `json:"commercialZone"`
	Warehouse      string `json:"warehouse"`
	RoleID         uint   `json:"roleId"`
	Password       string `json:"password"` // Opcional para actualizar
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Actualizar solo los campos proporcionados
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.CommercialZone != "" {
		user.CommercialZone = req.CommercialZone
	}
	if req.Warehouse != "" {
		user.Warehouse = req.Warehouse
	}
	if req.RoleID != 0 {
		user.RoleID = req.RoleID
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
