package dto

type RegisterMerchantRequest struct {
	BusinessName string `json:"business_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=8"`
}

type LoginMerchantRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type MerchantAuthResponse struct {
	Token        string `json:"token"`
	MerchantID   string `json:"merchant_id"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	Status       string `json:"status"`
}
type GetMeResponse struct {
    ID                 string  `json:"id"`
    BusinessName       string  `json:"business_name"`
    Email              string  `json:"email"`
    Status             string  `json:"status"`
    SolanaWallet       *string `json:"solana_wallet"`
    PolygonWallet      *string `json:"polygon_wallet"`
    EthereumWallet     *string `json:"ethereum_wallet"`
    KycSubmittedAt     *string `json:"kyc_submitted_at"`
    VerifiedAt         *string `json:"verified_at"`
    TotalEarningsUSD   string  `json:"total_earnings_usd"`
    CreatedAt          string  `json:"created_at"`
    // Password hash and token are intentionally omitted for security
}
