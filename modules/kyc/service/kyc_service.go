package service

import (
	"fmt"
	"time"

	"go-backend/modules/kyc/dto"
	"go-backend/modules/kyc/repository"
)

type KycService struct {
	repo *repository.KycRepository
}

func NewKycService(repo *repository.KycRepository) *KycService {
	return &KycService{repo: repo}
}

func (s *KycService) SubmitDocument(merchantID string, req dto.SubmitKycDocumentRequest) error {
	docID := fmt.Sprintf("kycdoc_%d", time.Now().UnixNano())
	return s.repo.UpsertDocument(docID, merchantID, req.DocType, req.FileUrl)
}

func (s *KycService) GetStatus(merchantID string) (*dto.KycStatusResponse, error) {
	return s.repo.GetMerchantKycDetails(merchantID)
}
