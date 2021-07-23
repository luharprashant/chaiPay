package models

import "time"

type Payment struct {
	ID            string    `json:"id"`
	UserId        string    `json:"userId"`
	Amount        string    `json:"amount"`
	TransactionId string    `json:"transactionId"`
	CreatedAt     time.Time `json:"created_at"`
}
