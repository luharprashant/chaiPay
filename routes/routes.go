package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luharprashant/chaiPay/controllers"
)

func Routes(router *gin.Engine) {
	router.GET("/chaipay/api/v1/healthcheck", healthcheck)
	router.GET("/chaipay/api/v1/get_charges", controllers.GetAllCharges)
	router.POST("/chaipay/api/v1/create_charge", controllers.CreateCharge)
	router.POST("/chaipay/api/v1/capture_charge/:chargeId", controllers.CaptureCharge)
	router.POST("/chaipay/api/v1/create_refund/:chargeId", controllers.RefundCharge)
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
