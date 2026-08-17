package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-backend/modules/wallet/dto"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

// ============================================================
// CREATE WALLET
// ============================================================

func (r *WalletRepository) Create(
	merchantID string,
	req dto.CreateWalletRequest,
) (*dto.WalletResponse, error) {

	const query = `
		INSERT INTO merchant_wallets (
			merchant_id,
			currency,
			network,
			wallet_address,
			is_active
		)
		VALUES ($1, $2, $3, $4, TRUE)
		RETURNING
			id,
			merchant_id,
			currency,
			network,
			wallet_address,
			is_active,
			created_at,
			updated_at
	`

	var (
		res       dto.WalletResponse
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRow(
		query,
		merchantID,
		req.Currency,
		req.Network,
		req.WalletAddress,
	).Scan(
		&res.ID,
		&res.MerchantID,
		&res.Currency,
		&res.Network,
		&res.WalletAddress,
		&res.IsActive,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create wallet: %w",
			err,
		)
	}

	res.CreatedAt = createdAt.Format(time.RFC3339)
	res.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &res, nil
}

// ============================================================
// LIST WALLETS
// ============================================================

func (r *WalletRepository) ListByMerchantID(
	merchantID string,
) ([]dto.WalletResponse, error) {

	const query = `
		SELECT
			id,
			merchant_id,
			currency,
			network,
			wallet_address,
			is_active,
			created_at,
			updated_at
		FROM merchant_wallets
		WHERE merchant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		query,
		merchantID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to list wallets: %w",
			err,
		)
	}

	defer rows.Close()

	wallets := make(
		[]dto.WalletResponse,
		0,
	)

	for rows.Next() {
		var (
			res       dto.WalletResponse
			createdAt time.Time
			updatedAt time.Time
		)

		if err := rows.Scan(
			&res.ID,
			&res.MerchantID,
			&res.Currency,
			&res.Network,
			&res.WalletAddress,
			&res.IsActive,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to scan wallet: %w",
				err,
			)
		}

		res.CreatedAt = createdAt.Format(time.RFC3339)
		res.UpdatedAt = updatedAt.Format(time.RFC3339)

		wallets = append(wallets, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed to iterate wallets: %w",
			err,
		)
	}

	return wallets, nil
}

// ============================================================
// GET WALLET
// ============================================================

func (r *WalletRepository) GetByID(
	merchantID string,
	walletID int64,
) (*dto.WalletResponse, error) {

	const query = `
		SELECT
			id,
			merchant_id,
			currency,
			network,
			wallet_address,
			is_active,
			created_at,
			updated_at
		FROM merchant_wallets
		WHERE id = $1
		  AND merchant_id = $2
	`

	var (
		res       dto.WalletResponse
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRow(
		query,
		walletID,
		merchantID,
	).Scan(
		&res.ID,
		&res.MerchantID,
		&res.Currency,
		&res.Network,
		&res.WalletAddress,
		&res.IsActive,
		&createdAt,
		&updatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("wallet not found")
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get wallet: %w",
			err,
		)
	}

	res.CreatedAt = createdAt.Format(time.RFC3339)
	res.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &res, nil
}

// ============================================================
// UPDATE WALLET
// ============================================================

func (r *WalletRepository) Update(
	merchantID string,
	walletID int64,
	req dto.UpdateWalletRequest,
) (*dto.WalletResponse, error) {

	const query = `
		UPDATE merchant_wallets
		SET
			currency = COALESCE(
				NULLIF($1, ''),
				currency
			),
			network = COALESCE(
				NULLIF($2, ''),
				network
			),
			wallet_address = COALESCE(
				NULLIF($3, ''),
				wallet_address
			),
			is_active = COALESCE(
				$4,
				is_active
			),
			updated_at = NOW()
		WHERE id = $5
		  AND merchant_id = $6
		RETURNING
			id,
			merchant_id,
			currency,
			network,
			wallet_address,
			is_active,
			created_at,
			updated_at
	`

	var (
		res       dto.WalletResponse
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRow(
		query,
		req.Currency,
		req.Network,
		req.WalletAddress,
		req.IsActive,
		walletID,
		merchantID,
	).Scan(
		&res.ID,
		&res.MerchantID,
		&res.Currency,
		&res.Network,
		&res.WalletAddress,
		&res.IsActive,
		&createdAt,
		&updatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New(
			"wallet not found or unauthorized",
		)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to update wallet: %w",
			err,
		)
	}

	res.CreatedAt = createdAt.Format(time.RFC3339)
	res.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &res, nil
}

// ============================================================
// DELETE WALLET
// ============================================================

func (r *WalletRepository) Delete(
	merchantID string,
	walletID int64,
) error {

	const query = `
		DELETE FROM merchant_wallets
		WHERE id = $1
		  AND merchant_id = $2
	`

	result, err := r.db.Exec(
		query,
		walletID,
		merchantID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to delete wallet: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"failed to check deleted wallet: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return errors.New("wallet not found")
	}

	return nil
}
