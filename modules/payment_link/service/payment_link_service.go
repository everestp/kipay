package service

import (
    "fmt"
    "time"

    "go-backend/modules/payment_link/dto"
    "go-backend/modules/payment_link/repository"
)

type PaymentLinkService struct {
    repo *repository.PaymentLinkRepository
}

func NewPaymentLinkService(repo *repository.PaymentLinkRepository) *PaymentLinkService {
    return &PaymentLinkService{repo: repo}
}

func (s *PaymentLinkService) CreatePaymentLink(merchantID string, req dto.CreatePaymentLinkRequest) (*dto.PaymentLinkResponse, error) {
    linkID := fmt.Sprintf("plink_%d", time.Now().UnixNano())
    return s.repo.Create(linkID, merchantID, req)
}

func (s *PaymentLinkService) GetPaymentLink(id string) (*dto.PaymentLinkResponse, error) {
    return s.repo.GetByID(id)
}

// Added to support fetching all payment links belonging to a specific merchant from context
func (s *PaymentLinkService) GetAllPaymentLinksByMerchant(merchantID string) ([]dto.PaymentLinkResponse, error) {
    return s.repo.GetAllByMerchant(merchantID)
}

// UpdatePaymentLink handles updating an existing payment link with ownership validation
func (s *PaymentLinkService) UpdatePaymentLink(paymentLinkID string, merchantID string, req dto.UpdatePaymentLinkRequest) (*dto.PaymentLinkResponse, error) {
    return s.repo.Update(paymentLinkID, merchantID, req)
}
