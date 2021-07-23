package controllers

import (
	"github.com/luharprashant/chaiPay/models"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

// DATABASE INSTANCE
var collection *mongo.Collection

func PaymentCollection(c *mongo.Database) {
	collection = c.Collection("payments")
}

func GetAllPayments(c *gin.Context) {
	var todos []models.Payment
	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	// Iterate through the returned cursor.
	for cursor.Next(context.TODO()) {
		var todo models.Payment
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Todos",
		"data":    todos,
	})
	return
}
