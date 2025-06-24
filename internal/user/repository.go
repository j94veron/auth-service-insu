package user

import (
	"errors"
	"github.com/j94veron/auth-service-insu/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	List() ([]models.User, error)

	// New feature to check user restrictions
	IsUserRestricted(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Role.Permissions").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found\n")
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Role.Permissions").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("mail not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Role.Permissions").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *repository) List() ([]models.User, error) {
	var users []models.User
	if err := r.db.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) IsUserRestricted(id uint) (bool, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return false, err
	}
	return false, nil
}
