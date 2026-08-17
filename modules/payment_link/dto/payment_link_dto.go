
package dto

type CreatePaymentLinkRequest struct {
    Title               string   `json:"title"`
    Description         string   `json:"description"`
    ImageUrl            string   `json:"image_url"`
    PricingType         string   `json:"pricing_type"` // e.g., "fixed" or "custom"
    AmountUSD           float64  `json:"amount_usd"`
    MinAmountUSD        float64  `json:"min_amount_usd"`
    PrimaryCurrency     string   `json:"primary_currency"`
    SupportedCurrencies []string `json:"supported_currencies"`
    SupportedNetworks   []string `json:"supported_networks"`
    SuccessMessage      string   `json:"success_message"`
    RedirectUrl         string   `json:"redirect_url"`
    ContinueButtonText  string   `json:"continue_button_text"`
    Color               string   `json:"color"`
}

// UpdatePaymentLinkRequest uses pointers so that only provided fields are updated.
type UpdatePaymentLinkRequest struct {
    Title               *string  `json:"title"`
    Description         *string  `json:"description"`
    ImageUrl            *string  `json:"image_url"`
    PricingType         *string  `json:"pricing_type"`
    AmountUSD           *float64 `json:"amount_usd"`
    MinAmountUSD        *float64 `json:"min_amount_usd"`
    SupportedCurrencies []string `json:"supported_currencies"`
    SupportedNetworks   []string `json:"supported_networks"`
    SuccessMessage      *string  `json:"success_message"`
    RedirectUrl         *string  `json:"redirect_url"`
    ContinueButtonText  *string  `json:"continue_button_text"`
    Color               *string  `json:"color"`
    IsActive            *bool    `json:"is_active"`
}

// PaymentLinkResponse represents the standard response object returned to the client.
type PaymentLinkResponse struct {
    ID                  string   `json:"id"`
    MerchantID          string   `json:"merchant_id"`
    Title               string   `json:"title"`
    Description         string   `json:"description"`
    ImageUrl            string   `json:"image_url"`
    PricingType         string   `json:"pricing_type"`
    AmountUSD           float64  `json:"amount_usd"`
    MinAmountUSD        float64  `json:"min_amount_usd"`
    SupportedCurrencies []string `json:"supported_currencies"`
    SupportedNetworks   []string `json:"supported_networks"`
    SuccessMessage      string   `json:"success_message"`
    RedirectUrl         string   `json:"redirect_url"`
    ContinueButtonText  string   `json:"continue_button_text"`
    TotalRevenueUSD     float64  `json:"total_revenue_usd"`
    Color               string   `json:"color"`
    IsActive            bool     `json:"is_active"`
    CreatedAt           string   `json:"created_at"`
}
