package models

type ChargeRequest struct {
	Amount        int64    `json:"amount"`
}

type ChargeRefundResponse struct {
	ID string `json:"id"`
	Error string `json:"error"`
}

type ChargeDatabase struct {
	ID            string    `json:"id"`
	Amount        int64    `json:"amount"`
	CreatedAt int64 `json:"created_at"`
	Captured bool `json:"captured"`
	Refunded bool `json:"refunded"`
	RefundId string `json:"refund_id"`
}
