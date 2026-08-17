package service

import (
	"go-backend/modules/wallet/dto"
	"go-backend/modules/wallet/repository"
)

type WalletService struct {
	repo *repository.WalletRepository
}

func NewWalletService(repo *repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) CreateWallet(
	merchantID string,
	req dto.CreateWalletRequest,
) (*dto.WalletResponse, error) {
	return s.repo.Create(merchantID, req)
}

func (s *WalletService) ListMerchantWallets(
	merchantID string,
) ([]dto.WalletResponse, error) {
	return s.repo.ListByMerchantID(merchantID)
}

func (s *WalletService) GetWalletByID(
	merchantID string,
	walletID int64,
) (*dto.WalletResponse, error) {
	return s.repo.GetByID(merchantID, walletID)
}

func (s *WalletService) UpdateWallet(
	merchantID string,
	walletID int64,
	req dto.UpdateWalletRequest,
) (*dto.WalletResponse, error) {
	return s.repo.Update(merchantID, walletID, req)
}

func (s *WalletService) DeleteWallet(
	merchantID string,
	walletID int64,
) error {
	return s.repo.Delete(merchantID, walletID)
}
