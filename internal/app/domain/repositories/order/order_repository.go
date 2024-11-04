package order

import (
	"errors"
	"github.com/google/uuid"
	"github.com/nade-harlow/ecom-api/internal/adapter/database"
	"github.com/nade-harlow/ecom-api/internal/app/domain/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	GetUserOrderByID(id, userID uuid.UUID) (*models.Order, error)
	GetOrdersByUserID(userID uuid.UUID) ([]models.Order, error)
	GetAllOrders() ([]models.Order, error)
	UpdateOrderStatus(orderID uuid.UUID, status models.OrderStatus) error
	CancelOrder(orderID uuid.UUID) error
}

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository() OrderRepository {
	return &OrderRepositoryImpl{db: database.GetDbConnection()}
}

func (r *OrderRepositoryImpl) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepositoryImpl) GetUserOrderByID(id, userID uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").First(&order, "id = ? AND user_id = ?", id, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

func (r *OrderRepositoryImpl) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").First(&order, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

func (r *OrderRepositoryImpl) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Find(&orders).Error
	return orders, err
}

func (r *OrderRepositoryImpl) GetOrdersByUserID(userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return orders, err
}

func (r *OrderRepositoryImpl) UpdateOrderStatus(orderID uuid.UUID, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *OrderRepositoryImpl) CancelOrder(orderID uuid.UUID) error {
	return r.db.Model(&models.Order{}).Where("id = ? AND status = ?", orderID, "Pending").Update("status", "Canceled").Error
}
