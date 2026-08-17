
package dto

type SubmitKycDocumentRequest struct {
	DocType string `json:"doc_type" validate:"required"` // BUSINESS_REGISTRATION, IDENTITY_PROOF, ADDRESS_PROOF, TAX_DOCUMENT
	FileUrl string `json:"file_url" validate:"required"` // S3 or secure storage link
}

type KycStatusResponse struct {
	MerchantID     string       `json:"merchant_id"`
	MerchantStatus string       `json:"merchant_status"` // PENDING_KYC, IN_REVIEW, VERIFIED, SUSPENDED, BLOCKED
	KycSubmittedAt *string      `json:"kyc_submitted_at,omitempty"`
	VerifiedAt     *string      `json:"verified_at,omitempty"`
	Documents      []KycDocItem `json:"documents"`
}

type KycDocItem struct {
	ID             string  `json:"id"`
	DocType        string  `json:"doc_type"`
	FileUrl        string  `json:"file_url"`
	Status         string  `json:"status"` // PENDING, APPROVED, REJECTED
	RejectionReason *string `json:"rejection_reason,omitempty"`
	CreatedAt      string  `json:"created_at"`
}
