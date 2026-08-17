package dto

type CreatePayoutRequest struct {
	AmountUSD          float64 `json:"amount_usd" validate:"required,gt=0"`
	Currency           string  `json:"currency" validate:"required"`           // e.g., USDT, USDC, SOL
	DestinationAddress string  `json:"destination_address" validate:"required"` // External wallet address
	Network            string  `json:"network" validate:"required"`            // solana, polygon, ethereum
}

type PayoutResponse struct {
	ID                 string  `json:"id"`
	MerchantID         string  `json:"merchant_id"`
	AmountUSD          float64 `json:"amount_usd"`
	Currency           string  `json:"currency"`
	DestinationAddress string  `json:"destination_address"`
	Network            string  `json:"network"`
	Status             string  `json:"status"` // PENDING, PROCESSING, COMPLETED, FAILED
	CreatedAt          string  `json:"created_at"`
}
