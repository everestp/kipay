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
