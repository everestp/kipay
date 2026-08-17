package repository

import (
	"database/sql"
	"fmt"
	"time"

	"go-backend/modules/admin/dto"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) ListMerchants(statusFilter string) ([]dto.AdminMerchantSummaryResponse, error) {
	query := `SELECT id, business_name, email, status, total_earnings_usd, created_at FROM merchants`
	var rows *sql.Rows
	var err error

	if statusFilter != "" {
		query += ` WHERE status = $1 ORDER BY created_at DESC`
		rows, err = r.db.Query(query, statusFilter)
	} else {
		query += ` ORDER BY created_at DESC`
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []dto.AdminMerchantSummaryResponse
	for rows.Next() {
		var m dto.AdminMerchantSummaryResponse
		var cAt time.Time
		if err := rows.Scan(&m.ID, &m.BusinessName, &m.Email, &m.Status, &m.TotalEarningsUSD, &cAt); err == nil {
			m.CreatedAt = cAt.Format(time.RFC3339)
			list = append(list, m)
		}
	}
	return list, nil
}

func (r *AdminRepository) UpdateMerchantStatusWithNotification(merchantID string, newStatus string, reason string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Update merchant status
	var updateQuery string
	if newStatus == "VERIFIED" {
		updateQuery = `UPDATE merchants SET status = $1, verified_at = NOW() WHERE id = $2`
	} else {
		updateQuery = `UPDATE merchants SET status = $1 WHERE id = $2`
	}
	_, err = tx.Exec(updateQuery, newStatus, merchantID)
	if err != nil {
		return err
	}

	// 2. Dispatch automated system notification message to merchant profile
	notifID := fmt.Sprintf("notif_%d", time.Now().UnixNano())
	title := fmt.Sprintf("Account Status Update: %s", newStatus)
	message := fmt.Sprintf("Your Pinecone merchant status has been updated to %s. Reason / Notes: %s", newStatus, reason)
	notifType := fmt.Sprintf("ACCOUNT_%s", newStatus)

	_, err = tx.Exec(
		`INSERT INTO merchant_notifications (id, merchant_id, title, message, type, is_read, created_at) VALUES ($1, $2, $3, $4, $5, FALSE, NOW())`,
		notifID, merchantID, title, message, notifType,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AdminRepository) ReviewKycDocument(docID string, adminID string, status string, rejectionReason *string) error {
	query := `UPDATE merchant_kyc_documents SET status = $1, reviewed_by = $2, reviewed_at = NOW(), rejection_reason = $3 WHERE id = $4`
	_, err := r.db.Exec(query, status, adminID, rejectionReason, docID)
	return err
}
