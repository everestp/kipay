package dto

type CreateApiKeyRequest struct {
	Name string `json:"name" validate:"required"` // e.g. "Production Storefront Key"
}

type ApiKeyResponse struct {
	ID         string `json:"id"`
	MerchantID string `json:"merchant_id"`
	Name       string `json:"name"`
	KeyPrefix  string `json:"key_prefix"`
	SecretKey  string `json:"secret_key,omitempty"` // Only returned once upon creation!
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
}
