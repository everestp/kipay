package service

import (
	"fmt"
	"time"

	"go-backend/modules/payout/dto"
	"go-backend/modules/payout/repository"
)

type PayoutService struct {
	repo *repository.PayoutRepository
}

func NewPayoutService(repo *repository.PayoutRepository) *PayoutService {
	return &PayoutService{repo: repo}
}

func (s *PayoutService) RequestPayout(merchantID string, req dto.CreatePayoutRequest) (*dto.PayoutResponse, error) {
	payoutID := fmt.Sprintf("po_%d", time.Now().UnixNano())
	return s.repo.CreatePayout(payoutID, merchantID, req)
}
