package service

import (
	"go-backend/modules/transaction/dto"
	"go-backend/modules/transaction/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetMerchantInvoices(merchantID string, status string, limit int, offset int) ([]dto.InvoiceHistoryItem, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListDirectInvoices(merchantID, status, limit, offset)
}
func (s *TransactionService) GetMerchantLinkInvoices(merchantID string , status string, limit int, offset int) ([]dto.LinkInvoiceHistoryItem, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListLinkInvoices(merchantID, status, limit, offset)
}
func (s *TransactionService) GetMerchantLinkInvoicesByPaymentLinkId(merchantID string, paymentlink_id string , status string, limit int, offset int) ([]dto.LinkInvoiceHistoryItem, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListLinkInvoicesByLinkId(merchantID, paymentlink_id,status, limit, offset)
}


func (s *TransactionService) GetMerchantSettledTransactions(merchantID string, limit int, offset int) ([]dto.TransactionHistoryItem, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListSettledTransactions(merchantID, limit, offset)
}
