package product

import (
	"errors"
	"github.com/google/uuid"
	"github.com/nade-harlow/ecom-api/internal/adapter/database"
	"github.com/nade-harlow/ecom-api/internal/app/domain/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id uuid.UUID) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	UpdateProduct(id uuid.UUID, update *models.Product) error
	DeleteProduct(id uuid.UUID) error
}

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{db: database.GetDbConnection()}
}

func (r *ProductRepositoryImpl) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepositoryImpl) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

func (r *ProductRepositoryImpl) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *ProductRepositoryImpl) UpdateProduct(id uuid.UUID, update *models.Product) error {
	return r.db.Where("id = ?", id).Save(update).Error
}

func (r *ProductRepositoryImpl) DeleteProduct(id uuid.UUID) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}
