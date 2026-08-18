package service

import (
	"go-backend/modules/dashboard/dto"
	"go-backend/modules/dashboard/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetMetrics(merchantID string) (*dto.DashboardMetricsResponse, error) {
	return s.repo.GetDashboardMetrics(merchantID)
}
