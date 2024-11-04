package order

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/dto"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/response"
	"github.com/nade-harlow/ecom-api/internal/app/domain/models"
	"github.com/nade-harlow/ecom-api/internal/app/domain/services/order"
	"github.com/nade-harlow/ecom-api/internal/app/utils/apperrors"
	utils "github.com/nade-harlow/ecom-api/internal/app/utils/auth"
	"github.com/nade-harlow/ecom-api/internal/app/utils/helper"
)

type OrderHandler struct {
	orderService order.OrderService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: order.NewOrderService(),
	}
}

func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderRequest
	var user utils.DecodedUser
	if err := ctx.ShouldBind(&req); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(helper.ValidatorFormatErrors(err).Error()))
		return
	}

	if err := helper.ValidateRequestBody(req); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(err.Error()))
		return
	}

	userC, found := ctx.Get("user")
	if !found {
		response.JsonError(ctx, apperrors.UnauthorizedError("Unauthorized. Login and try again"))
		return
	}

	mapstructure.Decode(userC, &user)
	userID, errS := helper.StringToUUID(user.UserID)
	if errS != nil {
		response.JsonError(ctx, apperrors.BadRequestError("invalid id format"))
		return
	}

	order, err := o.orderService.CreateOrder(userID, req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	response.JsonCreated(ctx, order, "order")
}

func (o *OrderHandler) CancelOrder(ctx *gin.Context) {
	paramID := ctx.Param("orderID")
	var user utils.DecodedUser

	if paramID == "" {
		response.JsonError(ctx, apperrors.BadRequestError("invalid order id"))
		return
	}

	orderID, errS := helper.StringToUUID(paramID)
	if errS != nil {
		response.JsonError(ctx, apperrors.BadRequestError("invalid id format"))
		return
	}

	userC, found := ctx.Get("user")
	if !found {
		response.JsonError(ctx, apperrors.UnauthorizedError("Unauthorized. Login and try again"))
		return
	}

	mapstructure.Decode(userC, &user)
	userID, errS := helper.StringToUUID(user.UserID)
	if errS != nil {
		response.JsonError(ctx, apperrors.BadRequestError("invalid id format"))
		return
	}

	err := o.orderService.CancelOrder(orderID, userID)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	response.JsonModified(ctx, nil, "order")
}

func (o *OrderHandler) UpdateOrderStatus(ctx *gin.Context) {
	paramID := ctx.Param("orderID")
	statusQ := ctx.Query("status")

	if paramID == "" {
		response.JsonError(ctx, apperrors.BadRequestError("invalid order id"))
		return
	}

	if statusQ == "" {
		response.JsonError(ctx, apperrors.BadRequestError("order status is required"))
		return
	}

	orderID, errS := helper.StringToUUID(paramID)
	if errS != nil {
		response.JsonError(ctx, apperrors.BadRequestError("invalid id format"))
		return
	}

	switch models.OrderStatus(statusQ) {
	case models.OrderStatus_Canceled, models.OrderStatus_Pending, models.OrderStatus_Fulfilled:
	default:
		response.JsonError(ctx, apperrors.BadRequestError("invalid order status"))

	}

	order, err := o.orderService.UpdateOrderStatus(orderID, models.OrderStatus(statusQ))
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	response.JsonModified(ctx, order, "order")
}

func (o *OrderHandler) GetOrderByID(ctx *gin.Context) {
	paramID := ctx.Param("orderID")

	if paramID == "" {
		response.JsonError(ctx, apperrors.BadRequestError("invalid order id"))
		return
	}

	orderID, errS := helper.StringToUUID(paramID)
	if errS != nil {
		response.JsonError(ctx, apperrors.BadRequestError("invalid id format"))
		return
	}

	order, err := o.orderService.GetOrderByID(orderID)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	response.JsonOk(ctx, order)
}

func (o *OrderHandler) GetUserOrders(ctx *gin.Context) {
	var user utils.DecodedUser

	userC, found := ctx.Get("user")
	if !found {
		response.JsonError(ctx, apperrors.UnauthorizedError("Unauthorized. Login and try again"))
		return
	}

	mapstructure.Decode(userC, &user)
	userID, errS := helper.StringToUUID(user.UserID)
	if errS != nil {
		response.JsonError(ctx, apperrors.BadRequestError("invalid id format"))
		return
	}

	orders, err := o.orderService.GetUserOrders(userID)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	response.JsonOk(ctx, orders)
}

func (o *OrderHandler) GetAllOrders(ctx *gin.Context) {

	orders, err := o.orderService.GetAllOrders()
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	response.JsonOk(ctx, orders)
}
