package repository

import (
    "database/sql"
    "errors"
    "fmt"
    "time"

    mo "go-backend/models"
)

type MerchantRepository struct {
    db *sql.DB
}

func NewMerchantRepository(db *sql.DB) *MerchantRepository {
    return &MerchantRepository{db: db}
}

func (r *MerchantRepository) Create(id string, businessName string, email string, passwordHash string) (*mo.Merchant, error) {
    query := `
        INSERT INTO merchants (id, business_name, email, password_hash, status, created_at)
        VALUES ($1, $2, $3, $4, 'PENDING_KYC', $5)
        RETURNING id, business_name, email, status, total_earnings_usd, created_at
    `
    var m mo.Merchant
    err := r.db.QueryRow(
        query,
        id, businessName, email, passwordHash, time.Now(),
    ).Scan(&m.ID, &m.BusinessName, &m.Email, &m.Status, &m.TotalEarningsUSD, &m.CreatedAt)

    if err != nil {
        return nil, err
    }
    return &m, nil
}

func (r *MerchantRepository) GetByEmail(email string) (*mo.Merchant, string, error) {
    query := `
        SELECT id, business_name, email, password_hash, status, total_earnings_usd, created_at
        FROM merchants
        WHERE email = $1
    `
    row := r.db.QueryRow(query, email)

    var m mo.Merchant
    var passwordHash string
    err := row.Scan(&m.ID, &m.BusinessName, &m.Email, &passwordHash, &m.Status, &m.TotalEarningsUSD, &m.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, "", errors.New("merchant not found")
    } else if err != nil {
        return nil, "", err
    }

    return &m, passwordHash, nil
}
func (r *MerchantRepository) GetByID(id string) (*mo.Merchant, string, error) {
    query := `
        SELECT id, business_name, email, password_hash, status, total_earnings_usd, created_at
        FROM merchants
        WHERE id = $1
    `
    row := r.db.QueryRow(query, id)

    var m mo.Merchant
    var passwordHash string
    err := row.Scan(&m.ID, &m.BusinessName, &m.Email, &passwordHash, &m.Status, &m.TotalEarningsUSD, &m.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, "", errors.New("merchant not found")
    } else if err != nil {
        return nil, "", err
    }

    return &m, passwordHash, nil
}

// ==========================================
// GET MERCHANT WALLET BY NETWORK (For Non-Custodial Payouts)
// ==========================================
func (r *MerchantRepository) GetWalletByNetwork(merchantID string, network string) (string, error) {
    var query string
    switch network {
    case "solana":
        query = `SELECT COALESCE(solana_wallet, '') FROM merchants WHERE id = $1`
    case "polygon":
        query = `SELECT COALESCE(polygon_wallet, '') FROM merchants WHERE id = $1`
    case "ethereum":
        query = `SELECT COALESCE(ethereum_wallet, '') FROM merchants WHERE id = $1`
    default:
        return "", fmt.Errorf("unsupported network: %s", network)
    }

    var wallet string
    err := r.db.QueryRow(query, merchantID).Scan(&wallet)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", errors.New("merchant not found")
        }
        return "", err
    }

    if wallet == "" {
        return "", fmt.Errorf("merchant has not configured a payout wallet for network %s", network)
    }

    return wallet, nil
}
