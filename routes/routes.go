package routes

import (
	"net/http"

	"github.com/luharprashant/chaiPay/controllers"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/api/v1/", healthcheck)
	router.GET("/api/v1/payments", controllers.GetAllPayments)
	router.NoRoute(notFound)
}

func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome, healthcheck complete",
	})
	return
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  404,
		"message": "Route Not Found",
	})
	return
}
