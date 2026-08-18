package repository

import (
	"database/sql"
	"fmt"

	"go-backend/modules/dashboard/dto"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}
func (r *DashboardRepository) GetDashboardMetrics(
	merchantID string,
) (*dto.DashboardMetricsResponse, error) {

	query := `
		SELECT
			COALESCE(
				SUM(amount_usd) FILTER (
					WHERE status = 'CONFIRMED'
				),
				0
			) AS total_volume_usd,

			COUNT(*) FILTER (
				WHERE status = 'CONFIRMED'
			) AS successful_payments_count,

			COUNT(*) FILTER (
				WHERE status = 'PENDING'
			) AS pending_invoices_count,

			COALESCE(
				(
					COUNT(*) FILTER (
						WHERE status = 'CONFIRMED'
					)::float
					/
					NULLIF(
						COUNT(*) FILTER (
							WHERE status IN ('CONFIRMED', 'FAILED', 'EXPIRED')
						),
						0
					)
				) * 100,
				0
			) AS success_rate

		FROM (
			SELECT
				amount_usd,
				status
			FROM payment_link_invoices
			WHERE merchant_id = $1

			UNION ALL

			SELECT
				amount_usd,
				status
			FROM direct_invoices
			WHERE merchant_id = $1
		) AS invoices
	`

	var metrics dto.DashboardMetricsResponse

	err := r.db.QueryRow(query, merchantID).Scan(
		&metrics.TotalVolumeUSD,
		&metrics.SuccessfulPaymentsCount,
		&metrics.PendingInvoicesCount,
		&metrics.SuccessRate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice dashboard metrics: %w", err)
	}

	// payment_links uses is_active, NOT status.
	activeLinksQuery := `
		SELECT COUNT(*)
		FROM payment_links
		WHERE merchant_id = $1
		  AND is_active = TRUE
	`

	err = r.db.QueryRow(activeLinksQuery, merchantID).
		Scan(&metrics.ActivePaymentLinksCount)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get active payment links count: %w",
			err,
		)
	}

	return &metrics, nil
}
