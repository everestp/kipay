package dto

// CreateLinkInvoiceRequest is used when a customer pays via a public payment link
type CreateLinkInvoiceRequest struct {
    PaymentLinkId string  `json:"payment_link_id" validate:"required"`
    AmountUSD     float64 `json:"amount_usd"` // Used only if the payment link allows custom amounts
    Currency      string  `json:"currency" validate:"required"` // USDT, USDC, POL, ETH
    Network       string  `json:"network" validate:"required"`  // polygon, solana, ethereum
}

// CreateDirectInvoiceRequest is used by external e-commerce plugins (WooCommerce, Shopify) via API Key
type CreateDirectInvoiceRequest struct {
    OrderID   string  `json:"order_id" validate:"required"`   // Cart/Order ID from the external store
    AmountUSD float64 `json:"amount_usd" validate:"required,gt=0"`
    Currency  string  `json:"currency" validate:"required"` // USDT, USDC, POL, ETH
    Network   string  `json:"network" validate:"required"`  // polygon, solana, ethereum
}

// InvoiceResponse is the unified response sent back to the checkout frontend or e-commerce API caller
type InvoiceResponse struct {
    InvoiceID      string  `json:"invoice_id"`
    OrderID       string   `json:"order_id"`
    MerchantID     string  `json:"merchant_id"`
    AmountUSD      float64 `json:"amount_usd"`
    Currency       string  `json:"currency"`
    Network        string  `json:"network"`
    AmountCrypto   float64 `json:"amount_crypto"`
    Status         string  `json:"status"` // PENDING, CONFIRMED, EXPIRED, FAILED
    DepositAddress string  `json:"deposit_address"`
    QRCodeData     string  `json:"qr_code_data"`
    ExpiresAt      string  `json:"expires_at"`
    ConfirmedAt    string  `json:"confirmed_at,omitempty"`
}
