package models

import (
    "time"
)

type Merchant struct {
    ID                string     `json:"id" db:"id"`
    BusinessName      string     `json:"business_name" db:"business_name"`
    Email             string     `json:"email" db:"email"`
    Status            string     `json:"status" db:"status"`
    SessionToken      *string    `json:"session_token,omitempty" db:"session_token"`
    SolanaWallet      *string    `json:"solana_wallet,omitempty" db:"solana_wallet"`
    PolygonWallet     *string    `json:"polygon_wallet,omitempty" db:"polygon_wallet"`
    EthereumWallet    *string    `json:"ethereum_wallet,omitempty" db:"ethereum_wallet"`
    KycSubmittedAt    *time.Time `json:"kyc_submitted_at,omitempty" db:"kyc_submitted_at"`
    VerifiedAt        *time.Time `json:"verified_at,omitempty" db:"verified_at"`
    TotalEarningsUSD  float64    `json:"total_earnings_usd" db:"total_earnings_usd"`
    CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}

type MerchantWallet struct {
    ID            string    `json:"id" db:"id"`
    MerchantID    string    `json:"merchant_id" db:"merchant_id"`
    Currency      string    `json:"currency" db:"currency"`
    Network       string    `json:"network" db:"network"`
    WalletAddress string    `json:"wallet_address" db:"wallet_address"`
    IsActive      bool      `json:"is_active" db:"is_active"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type AdminUser struct {
    ID           string    `json:"id" db:"id"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    Role         string    `json:"role" db:"role"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type APIKey struct {
    ID         string    `json:"id" db:"id"`
    MerchantID string    `json:"merchant_id" db:"merchant_id"`
    Name       string    `json:"name" db:"name"`
    KeyPrefix  string    `json:"key_prefix" db:"key_prefix"`
    SecretHash string    `json:"-" db:"secret_hash"`
    IsActive   bool      `json:"is_active" db:"is_active"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
}


type WebhookEndpointInfo struct {
	ID               string
	URL              string
	SubscribedEvents []string
	IsActive         bool
	CreatedAt        string
}

type MerchantKYCDocument struct {
    ID              string     `json:"id" db:"id"`
    MerchantID      string     `json:"merchant_id" db:"merchant_id"`
    DocType         string     `json:"doc_type" db:"doc_type"`
    FileURL         string     `json:"file_url" db:"file_url"`
    Status          string     `json:"status" db:"status"`
    RejectionReason *string    `json:"rejection_reason,omitempty" db:"rejection_reason"`
    ReviewedBy      *string    `json:"reviewed_by,omitempty" db:"reviewed_by"`
    ReviewedAt      *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
    CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

type MerchantNotification struct {
    ID         string    `json:"id" db:"id"`
    MerchantID string    `json:"merchant_id" db:"merchant_id"`
    Title      string    `json:"title" db:"title"`
    Message    string    `json:"message" db:"message"`
    Type       string    `json:"type" db:"type"`
    IsRead     bool      `json:"is_read" db:"is_read"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type PaymentLink struct {
    ID                  string    `json:"id" db:"id"`
    MerchantID          string    `json:"merchant_id" db:"merchant_id"`
    Title               string    `json:"title" db:"title"`
    Description         *string   `json:"description,omitempty" db:"description"`
    ImageURL            *string   `json:"image_url,omitempty" db:"image_url"`
    PricingType         string    `json:"pricing_type" db:"pricing_type"`
    AmountUSD           float64   `json:"amount_usd" db:"amount_usd"`
    MinAmountUSD        float64   `json:"min_amount_usd" db:"min_amount_usd"`
    SupportedCurrencies []string  `json:"supported_currencies" db:"supported_currencies"`
    SupportedNetworks   []string  `json:"supported_networks" db:"supported_networks"`
    SuccessMessage      *string   `json:"success_message,omitempty" db:"success_message"`
    RedirectURL         *string   `json:"redirect_url,omitempty" db:"redirect_url"`
    ContinueButtonText  *string   `json:"continue_button_text,omitempty" db:"continue_button_text"`
    TotalRevenueUSD     float64   `json:"total_revenue_usd" db:"total_revenue_usd"`
    IsActive            bool      `json:"is_active" db:"is_active"`
    CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

type Invoice struct {
    ID                  string     `json:"id" db:"id"`
    MerchantID          string     `json:"merchant_id" db:"merchant_id"`
    PaymentLinkID       *string    `json:"payment_link_id,omitempty" db:"payment_link_id"`
    OrderID             *string    `json:"order_id,omitempty" db:"order_id"`
    AmountUSD           float64    `json:"amount_usd" db:"amount_usd"`
    Currency            string     `json:"currency" db:"currency"`
    Network             string     `json:"network" db:"network"`
    AmountCrypto        float64    `json:"amount_crypto" db:"amount_crypto"`
    Status              string     `json:"status" db:"status"`
    DepositAddress      string     `json:"deposit_address" db:"deposit_address"`
    QRCodeData          string     `json:"qr_code_data" db:"qr_code_data"`
    ExpiresAt           time.Time  `json:"expires_at" db:"expires_at"`
    ConfirmedAt         *time.Time `json:"confirmed_at,omitempty" db:"confirmed_at"`
    CreatedAt           time.Time  `json:"created_at" db:"created_at"`
}

type Transaction struct {
    ID                    string    `json:"id" db:"id"`
    PaymentLinkInvoiceID *string    `json:"payment_link_invoice_id,omitempty" db:"payment_link_invoice_id"`
    DirectInvoiceID      *string    `json:"direct_invoice_id,omitempty" db:"direct_invoice_id"`
    MerchantID            string    `json:"merchant_id" db:"merchant_id"`
    TxHash                string    `json:"tx_hash" db:"tx_hash"`
    Network               string    `json:"network" db:"network"`
    AmountCrypto          float64   `json:"amount_crypto" db:"amount_crypto"`
    Currency              string    `json:"currency" db:"currency"`
    Status                string    `json:"status" db:"status"`
    BlockNumber           int64     `json:"block_number" db:"block_number"`
    CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

type BlacklistedAddress struct {
    ID        string    `json:"id" db:"id"`
    Address   string    `json:"address" db:"address"`
    Reason    *string   `json:"reason,omitempty" db:"reason"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type MerchantPayout struct {
    ID                 string    `json:"id" db:"id"`
    MerchantID         string    `json:"merchant_id" db:"merchant_id"`
    AmountUSD          float64   `json:"amount_usd" db:"amount_usd"`
    Currency           string    `json:"currency" db:"currency"`
    DestinationAddress string    `json:"destination_address" db:"destination_address"`
    Network            string    `json:"network" db:"network"`
    Status             string    `json:"status" db:"status"`
    CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type WebhookEndpoint struct {
    ID               string    `json:"id" db:"id"`
    MerchantID       string    `json:"merchant_id" db:"merchant_id"`
    EndpointURL      string    `json:"endpoint_url" db:"endpoint_url"`
    SubscribedEvents []string  `json:"subscribed_events" db:"subscribed_events"`
    Secret           string    `json:"secret" db:"secret"`
    IsActive         bool      `json:"is_active" db:"is_active"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type WebhookLog struct {
    ID           string    `json:"id" db:"id"`
    MerchantID   string    `json:"merchant_id" db:"merchant_id"`
    EventType    string    `json:"event_type" db:"event_type"`
    EndpointURL  string    `json:"endpoint_url" db:"endpoint_url"`
    Payload      string    `json:"payload" db:"payload"`
    ResponseCode *int      `json:"response_code,omitempty" db:"response_code"`
    Status       string    `json:"status" db:"status"`
    RetryCount   int       `json:"retry_count" db:"retry_count"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
