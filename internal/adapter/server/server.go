package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/nade-harlow/ecom-api/internal/config"
	"log"
	"time"
)

func StartUpServer() {
	port := config.AppConfig.AppPort

	router := gin.New()
	router.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		gin.Logger(),
		gin.Recovery(),
		requestid.New(requestid.WithCustomHeaderStrKey("x-request-id")),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is alive",
		})
	})

	SetRoutes(router)

	log.Println("application now running on port: ", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting http server", err)
	}

}
