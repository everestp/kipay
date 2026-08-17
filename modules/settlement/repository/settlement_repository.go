package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type SettlementRepository struct {
	db *sql.DB
}

func NewSettlementRepository(db *sql.DB) *SettlementRepository {
	return &SettlementRepository{db: db}
}

// IsAddressBlacklisted checks if an address exists in the blacklisted_addresses table
func (r *SettlementRepository) IsAddressBlacklisted(ctx context.Context, address string) (bool, error) {
	var id string
	query := `SELECT id FROM blacklisted_addresses WHERE address = $1`
	err := r.db.QueryRowContext(ctx, query, address).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// IsTxProcessed checks if a transaction hash has already been recorded (Replay Protection)
func (r *SettlementRepository) IsTxProcessed(ctx context.Context, txHash string) (bool, error) {
	var id string
	query := `SELECT id FROM transactions WHERE tx_hash = $1`
	err := r.db.QueryRowContext(ctx, query, txHash).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// RecordSettlement performs an atomic database transaction: locks invoice, records transaction, updates status, and credits merchant earnings
func (r *SettlementRepository) RecordSettlement(
	ctx context.Context,
	invoiceID string,
	txHash string,
	sender string,
	recipient string,
	amountCrypto float64,
	currency string,
	network string,
	blockNumber int64,
) (string, error) {

	// =========================================================
	// BEGIN ATOMIC TRANSACTION
	// =========================================================

	tx, err := r.db.BeginTx(
		ctx,
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to begin database transaction: %w",
			err,
		)
	}

	defer tx.Rollback()

	// =========================================================
	// 1. LOCK AND FETCH INVOICE
	// =========================================================

	var (
		merchantID     string
		invoiceStatus  string
		amountUSD      float64
		depositAddress string
	)

	invoiceQuery := `
		SELECT
			merchant_id,
			status,
			amount_usd,
			deposit_address
		FROM payment_link_invoices
		WHERE id = $1
		FOR UPDATE
	`

	err = tx.QueryRowContext(
		ctx,
		invoiceQuery,
		invoiceID,
	).Scan(
		&merchantID,
		&invoiceStatus,
		&amountUSD,
		&depositAddress,
	)

	if err == sql.ErrNoRows {
		return "", errors.New("invoice not found")
	}

	if err != nil {
		return "", fmt.Errorf(
			"failed to query invoice: %w",
			err,
		)
	}

	// =========================================================
	// 2. VERIFY RECEIVER
	// =========================================================

	if depositAddress != recipient {
		return "", fmt.Errorf(
			"invalid receiver address: got %s, expected %s",
			recipient,
			depositAddress,
		)
	}

	// =========================================================
	// 3. INVOICE MUST BE PENDING
	// =========================================================

	if invoiceStatus != "PENDING" {
		return "", fmt.Errorf(
			"invoice is not in PENDING state: %s",
			invoiceStatus,
		)
	}

	// =========================================================
	// 4. REPLAY PROTECTION
	// =========================================================

	var existingTxID int64

	replayQuery := `
		SELECT id
		FROM transactions
		WHERE tx_hash = $1
		LIMIT 1
	`

	err = tx.QueryRowContext(
		ctx,
		replayQuery,
		txHash,
	).Scan(&existingTxID)

	if err == nil {
		return "", fmt.Errorf(
			"replay attack: transaction already recorded with id %d",
			existingTxID,
		)
	}

	if err != sql.ErrNoRows {
		return "", fmt.Errorf(
			"failed replay check: %w",
			err,
		)
	}

	// =========================================================
	// 5. INSERT TRANSACTION
	// =========================================================

	insertTxQuery := `
		INSERT INTO transactions (
			invoice_id,
			merchant_id,
			tx_hash,
			sender_address,
			receiver_address,
			network,
			amount_crypto,
			currency,
			status,
			block_number
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)
		RETURNING id
	`

	var transactionID int64

	err = tx.QueryRowContext(
		ctx,
		insertTxQuery,
		invoiceID,
		merchantID,
		txHash,
		sender,
		recipient,
		network,
		amountCrypto,
		currency,
		"CONFIRMED",
		blockNumber,
	).Scan(&transactionID)

	if err != nil {
		return "", fmt.Errorf(
			"failed to insert transaction: %w",
			err,
		)
	}

	// =========================================================
	// 6. UPDATE INVOICE
	// =========================================================

	updateInvoiceQuery := `
		UPDATE payment_link_invoices
		SET
			status = 'CONFIRMED',
			confirmed_at = NOW()
		WHERE id = $1
		  AND status = 'PENDING'
	`

	result, err := tx.ExecContext(
		ctx,
		updateInvoiceQuery,
		invoiceID,
	)

	if err != nil {
		return "", fmt.Errorf(
			"failed to update invoice: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf(
			"failed to verify invoice update: %w",
			err,
		)
	}

	if rowsAffected != 1 {
		return "", errors.New(
			"invoice was not updated",
		)
	}

	// =========================================================
	// 7. UPDATE MERCHANT EARNINGS
	// =========================================================

	updateMerchantQuery := `
		UPDATE merchants
		SET total_earnings_usd =
			total_earnings_usd + $1
		WHERE id = $2
	`

	result, err = tx.ExecContext(
		ctx,
		updateMerchantQuery,
		amountUSD,
		merchantID,
	)

	if err != nil {
		return "", fmt.Errorf(
			"failed to update merchant earnings: %w",
			err,
		)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf(
			"failed to verify merchant update: %w",
			err,
		)
	}

	if rowsAffected != 1 {
		return "", errors.New(
			"merchant not found while updating earnings",
		)
	}

	// =========================================================
	// 8. COMMIT
	// =========================================================

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf(
			"failed to commit settlement transaction: %w",
			err,
		)
	}

	return strconv.FormatInt(transactionID, 10), nil
}
