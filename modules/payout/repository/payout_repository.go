package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-backend/modules/payout/dto"
)

type PayoutRepository struct {
	db *sql.DB
}

func NewPayoutRepository(db *sql.DB) *PayoutRepository {
	return &PayoutRepository{db: db}
}

func (r *PayoutRepository) CreatePayout(payoutID string, merchantID string, req dto.CreatePayoutRequest) (*dto.PayoutResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Lock and check merchant earnings balance
	var totalEarnings float64
	err = tx.QueryRow(`SELECT total_earnings_usd FROM merchants WHERE id = $1 FOR UPDATE`, merchantID).Scan(&totalEarnings)
	if err == sql.ErrNoRows {
		return nil, errors.New("merchant not found")
	} else if err != nil {
		return nil, err
	}

	// Calculate total already paid out or pending payout
	var totalPendingOrCompleted float64
	_ = tx.QueryRow(
		`SELECT COALESCE(SUM(amount_usd), 0.00) FROM merchant_payouts WHERE merchant_id = $1 AND status IN ('PENDING', 'PROCESSING', 'COMPLETED')`,
		merchantID,
	).Scan(&totalPendingOrCompleted)

	availableBalance := totalEarnings - totalPendingOrCompleted
	if req.AmountUSD > availableBalance {
		return nil, fmt.Errorf("insufficient balance: requested $%.2f, available $%.2f", req.AmountUSD, availableBalance)
	}

	// 2. Insert payout request record
	query := `
		INSERT INTO merchant_payouts (id, merchant_id, amount_usd, currency, destination_address, network, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'PENDING', NOW())
		RETURNING id, merchant_id, amount_usd, currency, destination_address, network, status, created_at
	`

	var res dto.PayoutResponse
	var createdAtTime time.Time

	err = tx.QueryRow(
		query,
		payoutID, merchantID, req.AmountUSD, req.Currency, req.DestinationAddress, req.Network,
	).Scan(
		&res.ID, &res.MerchantID, &res.AmountUSD, &res.Currency,
		&res.DestinationAddress, &res.Network, &res.Status, &createdAtTime,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create payout request: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res.CreatedAt = createdAtTime.Format(time.RFC3339)
	return &res, nil
}
