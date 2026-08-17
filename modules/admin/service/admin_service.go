package service

import (
	"go-backend/modules/admin/dto"
	"go-backend/modules/admin/repository"
)

type AdminService struct {
	repo *repository.AdminRepository
}

func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetAllMerchants(status string) ([]dto.AdminMerchantSummaryResponse, error) {
	return s.repo.ListMerchants(status)
}

func (s *AdminService) ModifyMerchantStatus(merchantID string, req dto.UpdateMerchantStatusRequest) error {
	return s.repo.UpdateMerchantStatusWithNotification(merchantID, req.Status, req.Reason)
}

func (s *AdminService) ReviewDocument(docID string, adminID string, req dto.ReviewKycRequest) error {
	return s.repo.ReviewKycDocument(docID, adminID, req.Status, req.RejectionReason)
}
