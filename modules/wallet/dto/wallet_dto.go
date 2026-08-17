package dto

type CreateWalletRequest struct {
    Currency      string `json:"currency" validate:"required"`           // e.g., USDT, USDC
    Network       string `json:"network" validate:"required"`            // e.g., solana, polygon
    WalletAddress string `json:"wallet_address" validate:"required"` // External wallet address
}

type UpdateWalletRequest struct {
    Currency      string `json:"currency"`
    Network       string `json:"network"`
    WalletAddress string `json:"wallet_address"`
    IsActive      *bool  `json:"is_active"` // pointer to allow explicit false/true updates
}

type WalletResponse struct {
    ID            string `json:"id"`
    MerchantID    string `json:"merchant_id"`
    Currency      string `json:"currency"`
    Network       string `json:"network"`
    WalletAddress string `json:"wallet_address"`
    IsActive      bool   `json:"is_active"`
    CreatedAt     string `json:"created_at"`
    UpdatedAt     string `json:"updated_at"`
}
