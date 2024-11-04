package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/nade-harlow/ecom-api/internal/adapter/database"
	"github.com/nade-harlow/ecom-api/internal/app/domain/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{db: database.GetDbConnection()}
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepositoryImpl) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}
