package main

import (
	"log"

	gin "github.com/gin-gonic/gin"

	"github.com/luharprashant/chaiPay/config"
	"github.com/luharprashant/chaiPay/routes"
)

func main() {
	// Database
	config.Connect()

	// Init Router
	router := gin.Default()

	// Route Handlers / Endpoints
	routes.Routes(router)

	log.Fatal(router.Run(":50051"))
}
