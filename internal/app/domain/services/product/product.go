package product

import (
	"github.com/google/uuid"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/dto"
	"github.com/nade-harlow/ecom-api/internal/app/domain/models"
	"github.com/nade-harlow/ecom-api/internal/app/domain/repositories/product"
	"github.com/nade-harlow/ecom-api/internal/app/utils/apperrors"
)

type ProductService interface {
	CreateProduct(req dto.CreateProductRequest) (models.Product, error)
	UpdateProduct(productID uuid.UUID, req *dto.UpdateProductRequest) (*models.Product, apperrors.AppError)
	DeleteProduct(productID uuid.UUID) apperrors.AppError
	GetProductByID(productID uuid.UUID) (*models.Product, apperrors.AppError)
	GetProducts() ([]models.Product, apperrors.AppError)
}

type productService struct {
	repository product.ProductRepository
}

func NewProductService() ProductService {
	return &productService{
		repository: product.NewProductRepository(),
	}
}

func (p *productService) CreateProduct(req dto.CreateProductRequest) (models.Product, error) {
	var product models.Product
	product.NewProduct(req.Name, req.Description, req.Price, req.Stock)

	return product, p.repository.CreateProduct(&product)
}

func (p *productService) UpdateProduct(productID uuid.UUID, req *dto.UpdateProductRequest) (*models.Product, apperrors.AppError) {
	product, err := p.repository.GetProductByID(productID)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	if product == nil {
		return nil, apperrors.NotFoundError("product doesn't exist")
	}

	product.UpdateProduct(req.Name, req.Description, req.Price, req.Stock)
	err = p.repository.UpdateProduct(productID, product)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	product, err = p.repository.GetProductByID(productID)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	return product, nil
}

func (p *productService) DeleteProduct(productID uuid.UUID) apperrors.AppError {
	product, err := p.repository.GetProductByID(productID)
	if err != nil {
		return apperrors.InternalServerError("something went wrong")
	}

	if product == nil {
		return apperrors.NotFoundError("product doesn't exist")
	}

	if err = p.repository.DeleteProduct(productID); err != nil {
		return apperrors.InternalServerError("something went wrong")
	}

	return nil
}

func (p *productService) GetProductByID(productID uuid.UUID) (*models.Product, apperrors.AppError) {
	product, err := p.repository.GetProductByID(productID)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	if product == nil {
		return nil, apperrors.NotFoundError("product doesn't exist")
	}

	return product, nil
}

func (p *productService) GetProducts() ([]models.Product, apperrors.AppError) {
	products, err := p.repository.GetAllProducts()
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	return products, nil
}
