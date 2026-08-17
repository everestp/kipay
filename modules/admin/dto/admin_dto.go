package dto

type ReviewKycRequest struct {
	MerchantID     string  `json:"merchant_id" validate:"required"`
	Status         string  `json:"status" validate:"required"` // APPROVED or REJECTED
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

type UpdateMerchantStatusRequest struct {
	Status string `json:"status" validate:"required"` // VERIFIED, SUSPENDED, BLOCKED
	Reason string `json:"reason" validate:"required"` // Reason dispatched to merchant notifications
}

type AdminMerchantSummaryResponse struct {
	ID               string  `json:"id"`
	BusinessName     string  `json:"business_name"`
	Email            string  `json:"email"`
	Status           string  `json:"status"`
	TotalEarningsUSD float64 `json:"total_earnings_usd"`
	CreatedAt        string  `json:"created_at"`
}
