package role

import (
	"errors"
	"github.com/j94veron/auth-service-insu/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(id uint) (*models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uint) error
	List() ([]models.Role, error)
	CheckPermission(roleID uint, endpoint, method string) (bool, error)
	GetRoleName(roleID uint) (string, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("rol no encontrado")
		}
		return nil, err
	}
	return &role, nil
}

func (r *repository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *repository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *repository) List() ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *repository) CheckPermission(roleID uint, endpoint, method string) (bool, error) {
	var count int64

	err := r.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.endpoint = ? AND permissions.method = ?",
			roleID, endpoint, method).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) GetRoleName(roleID uint) (string, error) {
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("rol no encontrado")
		}
		return "", err // Propagamos otros errores
	}
	return role.Name, nil // Devuelve el nombre del rol
}
