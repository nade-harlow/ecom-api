package dto

import (
	"github.com/google/uuid"
	"github.com/nade-harlow/ecom-api/internal/app/domain/models"
)

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,dive"`
}

type CreateOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status models.OrderStatus `json:"status" binding:"required,oneof=Pending Fulfilled Canceled"`
}
