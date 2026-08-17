package dto

type ProcessTransactionRequest struct {
	InvoiceID      string  `json:"invoice_id" validate:"required"`
	TxHash         string  `json:"tx_hash" validate:"required"` // Blockchain transaction signature/hash
	Network        string  `json:"network" validate:"required"` // solana, polygon, ethereum
	AmountPaid     float64 `json:"amount_paid" validate:"required,gt=0"`
	Currency       string  `json:"currency" validate:"required"`
	FromAddress    string  `json:"from_address" validate:"required"`
	BlockNumber    int64   `json:"block_number"`
}
