package repository

import (
	"database/sql"
	"go-backend/modules/dashboard/dto"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetMerchantMetrics(merchantID string) (*dto.DashboardMetricsResponse, error) {
	var metrics dto.DashboardMetricsResponse

	// 1. Total volume and successful payments count from confirmed invoices / transactions
	queryVolume := `
		SELECT COALESCE(SUM(amount_usd), 0.00), COUNT(id)
		FROM invoices
		WHERE merchant_id = $1 AND status = 'CONFIRMED'
	`
	_ = r.db.QueryRow(queryVolume, merchantID).Scan(&metrics.TotalVolumeUSD, &metrics.SuccessfulPaymentsCount)

	// 2. Pending invoices count
	queryPending := `
		SELECT COUNT(id)
		FROM invoices
		WHERE merchant_id = $1 AND status = 'PENDING'
	`
	_ = r.db.QueryRow(queryPending, merchantID).Scan(&metrics.PendingInvoicesCount)

	// 3. Active payment links count
	queryLinks := `
		SELECT COUNT(id)
		FROM payment_links
		WHERE merchant_id = $1 AND is_active = TRUE
	`
	_ = r.db.QueryRow(queryLinks, merchantID).Scan(&metrics.ActivePaymentLinksCount)

	// 4. Calculate success rate percentage
	var totalInvoices int
	queryTotal := `SELECT COUNT(id) FROM invoices WHERE merchant_id = $1 AND status != 'PENDING'`
	_ = r.db.QueryRow(queryTotal, merchantID).Scan(&totalInvoices)

	if totalInvoices > 0 {
		metrics.SuccessRate = (float64(metrics.SuccessfulPaymentsCount) / float64(totalInvoices)) * 100.0
	} else {
		metrics.SuccessRate = 100.0
	}

	return &metrics, nil
}
