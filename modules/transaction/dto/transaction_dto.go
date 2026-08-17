package dto

type InvoiceHistoryItem struct {
	ID             string  `json:"id"`
	AmountUSD      float64 `json:"amount_usd"`
	Currency       string  `json:"currency"`
	Network        string  `json:"network"`
	AmountCrypto   float64 `json:"amount_crypto"`
	Status         string  `json:"status"` // PENDING, CONFIRMED, EXPIRED, FAILED
	DepositAddress string  `json:"deposit_address"`
	CreatedAt      string  `json:"created_at"`
	ConfirmedAt    *string `json:"confirmed_at,omitempty"`
}
type LinkInvoiceHistoryItem struct {
	ID             string  `json:"id"`
	AmountUSD      float64 `json:"amount_usd"`
	Currency       string  `json:"currency"`
	Network        string  `json:"network"`
	AmountCrypto   float64 `json:"amount_crypto"`
	Status         string  `json:"status"` // PENDING, CONFIRMED, EXPIRED, FAILED
	DepositAddress string  `json:"deposit_address"`
	PaymentLinkId    string     `json:"paymentlink_id"`
	CreatedAt      string  `json:"created_at"`
	ConfirmedAt    *string `json:"confirmed_at,omitempty"`
}

type TransactionHistoryItem struct {
	ID          string  `json:"id"`
	InvoiceID   string  `json:"invoice_id"`
	TxHash      string  `json:"tx_hash"`
	Network     string  `json:"network"`
	AmountCrypto float64 `json:"amount_crypto"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	BlockNumber int64   `json:"block_number"`
	CreatedAt   string  `json:"created_at"`
}
