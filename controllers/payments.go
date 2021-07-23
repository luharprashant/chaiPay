package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/luharprashant/chaiPay/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/refund"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

// DATABASE INSTANCE
var collection *mongo.Collection
var key = "sk_test_51JGHIASAauWGzIg3mBYsGFgOqK7cNIfykI31yPYQrlIxu0KgiKD6VXw1VeqbbLeQznvHUR3XkKgXUYs38FQpbnJQ00PPlFNvGW"

func PaymentCollection(c *mongo.Database) {
	collection = c.Collection("payments")
}

func CreateCharge(c *gin.Context) {
	var chargeReq models.ChargeRequest
	c.BindJSON(&chargeReq)
	amount := chargeReq.Amount

	stripe.Key = key

	customerParams := &stripe.CustomerParams{
		Name: stripe.String("Jenny Rosen"),
		Address: &stripe.AddressParams{
			Line1: stripe.String("510 Townsend St"),
			PostalCode: stripe.String("98140"),
			City: stripe.String("San Francisco"),
			State: stripe.String("CA"),
			Country: stripe.String("US"),
		},
	}
	customer, err := customer.New(customerParams)
	if err != nil{
		log.Printf("Error while creating a new user\n")
		code := stripeErrorLogger(err)
		resp := models.ChargeRefundResponse{Error: err.Error()}
		c.JSON(code, resp)
		return
	}

	cardParams := &stripe.CardParams{
		Customer: stripe.String(customer.ID),
		Token: stripe.String("tok_mastercard"),
	}
	custCard, err := card.New(cardParams)
	if err != nil{
		log.Printf("Error while creating a new card\n")
		code := stripeErrorLogger(err)
		resp := models.ChargeRefundResponse{Error: err.Error()}
		c.JSON(code, resp)
		return
	}
	var capture = false
	chargeParams := &stripe.ChargeParams{
		Amount: stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("test charge"),
		Source: &stripe.SourceParams{Token: stripe.String(custCard.ID)},
		Customer: stripe.String(customer.ID),
		Capture: &capture,
	}
	resp, err := charge.New(chargeParams)
	if err != nil {
		log.Printf("Error while creating a new charge\n")
		code := stripeErrorLogger(err)
		resp := models.ChargeRefundResponse{Error: err.Error()}
		c.JSON(code, resp)
		return
	}

	newCharge := models.ChargeDatabase{
		ID:        resp.ID,
		Amount:    resp.Amount,
		CreatedAt: resp.Created,
		Refunded: resp.Refunded,
		Captured: resp.Captured,
	}

	_, err = collection.InsertOne(context.TODO(), newCharge)

	if err != nil {
		log.Printf("Error while inserting new charge into db but was created with ID,%s\n", resp.ID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error while inserting new charge into db but was created with ID " + resp.ID,
		})
		return
	}

	res := models.ChargeRefundResponse{
		ID:    resp.ID,
	}

	c.JSON(http.StatusCreated, res)
	return

}

func CaptureCharge(c *gin.Context) {

	chargeId := c.Param("chargeId")

	stripe.Key = key
	captureParams := stripe.CaptureParams{
		Params:                    stripe.Params{},
		Amount:                    stripe.Int64(100),
	}
	_, err := charge.Capture(chargeId, &captureParams)
	if err != nil{
		log.Printf("Error while capturing the charge\n")
		code := stripeErrorLogger(err)
		resp := gin.H{"error": err.Error()}
		c.JSON(code, resp)
		return
	}

	newData := bson.M{
		"$set": bson.M{
			"captured":       true,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"id": chargeId}, newData)
	if err != nil {
		log.Printf("Error while updating in db but was updated at merchant site, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"message":  "Error while updating in db but was updated at merchant site",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Charge Captured Successfully",
	})
	return

}

func RefundCharge(c *gin.Context) {

	chargeId := c.Param("chargeId")

	stripe.Key = key
	refundParams := stripe.RefundParams{
		Charge: stripe.String(chargeId),
	}
	refund, err := refund.New(&refundParams)
	if err != nil{
		log.Printf("Error while refunding the charge\n")
		code := stripeErrorLogger(err)
		resp := gin.H{"error": err.Error()}
		c.JSON(code, resp)
		return
	}

	newData := bson.M{
		"$set": bson.M{
			"refunded":       true,
			"refund_id": refund.ID,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"id": chargeId}, newData)
	if err != nil {
		log.Printf("Error while updating in db but was updated at merchant site, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"message":  "Error while updating in db but was updated at merchant site",
		})
		return
	}

	res := models.ChargeRefundResponse{
		ID:    refund.ID,
	}

	c.JSON(http.StatusOK, res)
	return

}

func GetAllCharges(c *gin.Context) {
	var charges []models.ChargeDatabase
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
		var charge models.ChargeDatabase
		err := cursor.Decode(&charge)
		if err != nil {
			log.Printf("Error while decoding, Reason: %v\n", err)
			return
		}
		charges = append(charges, charge)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    charges,
	})
	return
}

func stripeErrorLogger(err error) int {

	if stripeErr, ok := err.(*stripe.Error); ok {
		// The Code field will contain a basic identifier for the failure.
		switch stripeErr.Code {
		case stripe.ErrorCodeCardDeclined:
		case stripe.ErrorCodeExpiredCard:
		case stripe.ErrorCodeIncorrectCVC:
		case stripe.ErrorCodeIncorrectZip:
			// etc.
		}

		// The Err field can be coerced to a more specific error type with a type
		// assertion. This technique can be used to get more specialized
		// information for certain errors.
		if cardErr, ok := stripeErr.Err.(*stripe.CardError); ok {
			log.Printf("Card was declined with code: %v\n", cardErr.DeclineCode)
		} else {
			log.Printf("Other Stripe error occurred: %v\n", stripeErr.Error())
		}
		return stripeErr.HTTPStatusCode
	} else {
		log.Printf("Other error occurred: %v\n", err.Error())
		return 500
	}

}