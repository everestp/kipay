package repository

import (
	"database/sql"
	"time"

	"go-backend/modules/transaction/dto"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) ListDirectInvoices(merchantID string, statusFilter string, limit int, offset int) ([]dto.InvoiceHistoryItem, error) {
	query := `
		SELECT id, amount_usd, currency, network, amount_crypto, status, deposit_address, created_at, confirmed_at
		FROM direct_invoices
		WHERE merchant_id = $1
	`
	var rows *sql.Rows
	var err error

	if statusFilter != "" {
		query += ` AND status = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		rows, err = r.db.Query(query, merchantID, statusFilter, limit, offset)
	} else {
		query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		rows, err = r.db.Query(query, merchantID, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []dto.InvoiceHistoryItem
	for rows.Next() {
		var item dto.InvoiceHistoryItem
		var cAt time.Time
		var confAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.AmountUSD, &item.Currency, &item.Network, &item.AmountCrypto, &item.Status, &item.DepositAddress, &cAt, &confAt); err == nil {
			item.CreatedAt = cAt.Format(time.RFC3339)
			if confAt.Valid {
				confStr := confAt.Time.Format(time.RFC3339)
				item.ConfirmedAt = &confStr
			}
			list = append(list, item)
		}
	}
	return list, nil
}
func (r *TransactionRepository) ListLinkInvoices(merchantID string, statusFilter string, limit int, offset int) ([]dto.LinkInvoiceHistoryItem, error) {
	query := `
		SELECT id, amount_usd, currency, network, payment_link_id, amount_crypto, status, deposit_address, created_at, confirmed_at
		FROM payment_link_invoices
		WHERE merchant_id = $1
	`
	var rows *sql.Rows
	var err error

	if statusFilter != "" {
		query += ` AND status = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		rows, err = r.db.Query(query, merchantID, statusFilter, limit, offset)
	} else {
		query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		rows, err = r.db.Query(query, merchantID, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []dto.LinkInvoiceHistoryItem
	for rows.Next() {
		var item dto.LinkInvoiceHistoryItem
		var cAt time.Time
		var confAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.AmountUSD, &item.Currency, &item.Network,&item.PaymentLinkId, &item.AmountCrypto, &item.Status, &item.DepositAddress, &cAt, &confAt); err == nil {
			item.CreatedAt = cAt.Format(time.RFC3339)
			if confAt.Valid {
				confStr := confAt.Time.Format(time.RFC3339)
				item.ConfirmedAt = &confStr
			}
			list = append(list, item)
		}
	}
	return list, nil
}

func (r *TransactionRepository) ListLinkInvoicesByLinkId(merchantID string, paymetlink_id string, statusFilter string, limit int, offset int) ([]dto.LinkInvoiceHistoryItem, error) {
  query := `
    SELECT id, amount_usd, currency, network, amount_crypto, status,
           deposit_address, payment_link_id, created_at, confirmed_at
    FROM payment_link_invoices
    WHERE merchant_id = $1 AND payment_link_id = $2
`

    var rows *sql.Rows
    var err error

    if statusFilter != "" {
        query += ` AND status = $3 ORDER BY created_at DESC LIMIT $4 OFFSET $5`
        rows, err = r.db.Query(query, merchantID, paymetlink_id, statusFilter, limit, offset)
    } else {
        query += ` ORDER BY created_at DESC LIMIT $3 OFFSET $4`
        rows, err = r.db.Query(query, merchantID, paymetlink_id, limit, offset)
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []dto.LinkInvoiceHistoryItem
    for rows.Next() {
        var item dto.LinkInvoiceHistoryItem
        var cAt time.Time
        var confAt sql.NullTime

        // Fixed: Scanned payment_link_id into item.PaymentLinkId instead of &item
        if err := rows.Scan(&item.ID, &item.AmountUSD, &item.Currency, &item.Network, &item.AmountCrypto, &item.Status, &item.DepositAddress, &item.PaymentLinkId, &cAt, &confAt); err == nil {
            item.CreatedAt = cAt.Format(time.RFC3339)
            if confAt.Valid {
                confStr := confAt.Time.Format(time.RFC3339)
                item.ConfirmedAt = &confStr
            }
            list = append(list, item)
        }
    }
    return list, nil
}

func (r *TransactionRepository) ListSettledTransactions(merchantID string, limit int, offset int) ([]dto.TransactionHistoryItem, error) {
	query := `
		SELECT id, invoice_id, tx_hash, network, amount_crypto, currency, status, block_number, created_at
		FROM transactions
		WHERE merchant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, merchantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []dto.TransactionHistoryItem
	for rows.Next() {
		var item dto.TransactionHistoryItem
		var cAt time.Time
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.TxHash, &item.Network, &item.AmountCrypto, &item.Currency, &item.Status, &item.BlockNumber, &cAt); err == nil {
			item.CreatedAt = cAt.Format(time.RFC3339)
			list = append(list, item)
		}
	}
	return list, nil
}
