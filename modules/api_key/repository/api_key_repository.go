package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-backend/modules/api_key/dto"
)

type ApiKeyRepository struct {
	db *sql.DB
}

func NewApiKeyRepository(db *sql.DB) *ApiKeyRepository {
	return &ApiKeyRepository{db: db}
}

func (r *ApiKeyRepository) Create(id string, merchantID string, name string, keyPrefix string, secretHash string) (*dto.ApiKeyResponse, error) {
	query := `
		INSERT INTO api_keys (id, merchant_id, name, key_prefix, secret_hash, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, TRUE, NOW())
		RETURNING id, merchant_id, name, key_prefix, is_active, created_at
	`
	var res dto.ApiKeyResponse
	var createdAt time.Time

	err := r.db.QueryRow(query, id, merchantID, name, keyPrefix, secretHash).Scan(
		&res.ID, &res.MerchantID, &res.Name, &res.KeyPrefix, &res.IsActive, &createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create api key: %v", err)
	}

	res.CreatedAt = createdAt.Format(time.RFC3339)
	return &res, nil
}

func (r *ApiKeyRepository) ValidatePrefixAndHash(keyPrefix string) (string, string, error) {
	query := `SELECT id, merchant_id, secret_hash, is_active FROM api_keys WHERE key_prefix = $1`
	var id, merchantID, secretHash string
	var isActive bool

	err := r.db.QueryRow(query, keyPrefix).Scan(&id, &merchantID, &secretHash, &isActive)
	if err == sql.ErrNoRows {
		return "", "", errors.New("invalid api key prefix")
	}
	if !isActive {
		return "", "", errors.New("api key is deactivated")
	}

	return merchantID, secretHash, nil
}
func (r *ApiKeyRepository) Delete(id string, merchantID string) error {
	query := `
		DELETE FROM api_keys
		WHERE id = $1 AND merchant_id = $2
	`

	result, err := r.db.Exec(query, id, merchantID)
	if err != nil {
		return fmt.Errorf("failed to delete api key: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check deleted api key: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("api key not found")
	}

	return nil
}
func (r *ApiKeyRepository) GetAllByMerchant(merchantID string) ([]dto.ApiKeyResponse, error) {
	query := `
		SELECT
			id,
			merchant_id,
			name,
			key_prefix,
			is_active,
			created_at
		FROM api_keys
		WHERE merchant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch api keys: %w", err)
	}
	defer rows.Close()

	var keys []dto.ApiKeyResponse

	for rows.Next() {
		var key dto.ApiKeyResponse
		var createdAt time.Time

		err := rows.Scan(
			&key.ID,
			&key.MerchantID,
			&key.Name,
			&key.KeyPrefix,
			&key.IsActive,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan api key: %w", err)
		}

		key.CreatedAt = createdAt.Format(time.RFC3339)

		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate api keys: %w", err)
	}

	return keys, nil
}
