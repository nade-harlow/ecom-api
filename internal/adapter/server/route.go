package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/handler/order"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/handler/product"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/handler/user"
)

const routePrefix = "/api/v1"

func SetRoutes(router *gin.Engine) {
	r := router.Group(routePrefix)
	userHandler := user.NewUserHandler()
	productHandler := product.NewProductHandler()
	orderHandler := order.NewOrderHandler()

	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/register", userHandler.RegisterUser)
		userRoutes.POST("/login", userHandler.Login)
	}

	productRoutes := r.Group("/product")
	{
		protectedRoute := productRoutes.Use(JwtMiddleware(), IsAdmin())
		protectedRoute.POST("/", productHandler.CreateProduct)
		protectedRoute.PATCH("/:productID", productHandler.UpdateProduct)
		protectedRoute.DELETE("/:productID", productHandler.DeleteProduct)
		productRoutes.GET("/:productID", JwtMiddleware(), productHandler.GetProductByID)
		productRoutes.GET("/", JwtMiddleware(), productHandler.GetProducts)
	}

	orderRoutes := r.Group("/order")
	{
		orderRoutes.POST("/", JwtMiddleware(), orderHandler.CreateOrder)
		orderRoutes.PATCH("/cancel/:orderID", JwtMiddleware(), orderHandler.CancelOrder)
		orderRoutes.PATCH("/:orderID", JwtMiddleware(), IsAdmin(), orderHandler.UpdateOrderStatus)
		orderRoutes.GET("/:orderID", JwtMiddleware(), orderHandler.GetOrderByID)
		orderRoutes.GET("/user", JwtMiddleware(), orderHandler.GetUserOrders)
		orderRoutes.GET("/", JwtMiddleware(), orderHandler.GetAllOrders)
	}
}
