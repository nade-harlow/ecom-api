package order

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/dto"
	"github.com/nade-harlow/ecom-api/internal/app/domain/models"
	"github.com/nade-harlow/ecom-api/internal/app/domain/repositories/order"
	"github.com/nade-harlow/ecom-api/internal/app/domain/repositories/product"
	"github.com/nade-harlow/ecom-api/internal/app/utils/apperrors"
	"log"
)

type OrderService interface {
	CreateOrder(userID uuid.UUID, req dto.CreateOrderRequest) (*models.Order, apperrors.AppError)
	CancelOrder(orderID, userID uuid.UUID) apperrors.AppError
	UpdateOrderStatus(orderID uuid.UUID, status models.OrderStatus) (*models.Order, apperrors.AppError)
	GetOrderByID(orderID uuid.UUID) (*models.Order, apperrors.AppError)
	GetUserOrders(userID uuid.UUID) ([]models.Order, apperrors.AppError)
	GetAllOrders() ([]models.Order, apperrors.AppError)
}

type orderService struct {
	repository        order.OrderRepository
	productRepository product.ProductRepository
}

func NewOrderService() OrderService {
	return &orderService{
		repository:        order.NewOrderRepository(),
		productRepository: product.NewProductRepository(),
	}
}

func (o *orderService) CreateOrder(userID uuid.UUID, req dto.CreateOrderRequest) (*models.Order, apperrors.AppError) {
	var order models.Order
	for _, item := range req.Items {
		product, err := o.productRepository.GetProductByID(item.ProductID)
		if err != nil {
			return nil, apperrors.InternalServerError("something went wrong")
		}

		if item.Quantity > product.Stock {
			fmt.Println("111111")
			return nil, apperrors.BadRequestError(fmt.Sprintf("%s quantity exceeds product stock", product.Name))
		}
		order.Items = append(order.Items, models.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

	}
	order.UserID = userID
	order.NewOrder()

	err := o.repository.CreateOrder(&order)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	//update product stock
	//todo: use transactions to ensure atomicity
	go func() {
		o.updateStock(order.Items, true)
	}()

	createdOrder, err := o.repository.GetOrderByID(order.ID)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	return createdOrder, nil
}

func (o *orderService) updateStock(orderItems []models.OrderItem, isRemove bool) {
	for _, orderItem := range orderItems {
		product, err := o.productRepository.GetProductByID(orderItem.ProductID)
		if err != nil {
			log.Println(fmt.Sprintf("failed to update product %s stock for order: %s", orderItem.ProductID, orderItem.OrderID))
			continue
		}

		if isRemove {
			product.Stock -= orderItem.Quantity
		} else {
			product.Stock += orderItem.Quantity
		}

		err = o.productRepository.UpdateProduct(product.ID, product)
		if err != nil {
			log.Println(fmt.Sprintf("failed to updatee product %s stock for order: %s", orderItem.ProductID, orderItem.OrderID))
			continue
		}
	}

	return
}

func (o *orderService) CancelOrder(orderID, userID uuid.UUID) apperrors.AppError {
	order, err := o.repository.GetUserOrderByID(orderID, userID)
	if err != nil {
		return apperrors.InternalServerError("something went wrong")
	}

	if order.Status != models.OrderStatus_Pending {
		return apperrors.BadRequestError(fmt.Sprintf("order has already been %s", order.Status))
	}

	err = o.repository.CancelOrder(orderID)
	if err != nil {
		return apperrors.InternalServerError("something went wrong")
	}

	//update product stock
	//todo: use transactions to ensure atomicity
	go func() {
		o.updateStock(order.Items, false)
	}()

	return nil
}

func (o *orderService) UpdateOrderStatus(orderID uuid.UUID, status models.OrderStatus) (*models.Order, apperrors.AppError) {
	order, err := o.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	if order == nil {
		return nil, apperrors.NotFoundError("order not found")
	}

	order.Status = status
	err = o.repository.UpdateOrderStatus(orderID, status)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	if status == models.OrderStatus_Canceled {
		//update product stock
		//todo: use transactions to ensure atomicity
		go func() {
			o.updateStock(order.Items, false)
		}()
	}

	return order, nil
}

func (o *orderService) GetOrderByID(orderID uuid.UUID) (*models.Order, apperrors.AppError) {
	order, err := o.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	if order == nil {
		return nil, apperrors.NotFoundError("order not found")
	}

	return order, nil
}

func (o *orderService) GetUserOrders(userID uuid.UUID) ([]models.Order, apperrors.AppError) {
	orders, err := o.repository.GetOrdersByUserID(userID)
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	return orders, nil
}

func (o *orderService) GetAllOrders() ([]models.Order, apperrors.AppError) {
	orders, err := o.repository.GetAllOrders()
	if err != nil {
		return nil, apperrors.InternalServerError("something went wrong")
	}

	return orders, nil
}
