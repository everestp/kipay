package repository

import (
    "database/sql"
    "fmt"
    "time"

    mo "go-backend/models"
)

type InvoiceRepository struct {
    db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
    return &InvoiceRepository{db: db}
}

// ==========================================
// 1. GET PAYMENT LINK BY ID
// ==========================================
type PaymentLinkEntity struct {
    ID          string
    MerchantID  string
    AmountUSD   float64
    MinAmountUSD float64
    PricingType string
    IsActive    bool
}

func (r *InvoiceRepository) GetPaymentLinkByID(paymentLinkId string) (*PaymentLinkEntity, error) {
    query := `
        SELECT id, merchant_id, amount_usd, min_amount_usd, pricing_type, is_active
        FROM payment_links
        WHERE id = $1
    `
    var link PaymentLinkEntity
    err := r.db.QueryRow(query, paymentLinkId).Scan(
        &link.ID, &link.MerchantID, &link.AmountUSD, &link.MinAmountUSD, &link.PricingType, &link.IsActive,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("payment link not found")
        }
        return nil, err
    }
    return &link, nil
}

// ==========================================
// 2. CREATE PAYMENT LINK INVOICE
// ==========================================
func (r *InvoiceRepository) CreatePaymentLinkInvoice(
    invoiceID string,
    merchantID string,
    paymentLinkId string,
    amountUSD float64,
    currency string,
    network string,
    amountCrypto float64,
    depositAddress string,
    qrCodeData string,
    expiresAt time.Time,
) (*mo.Invoice, error) {
    query := `
        INSERT INTO payment_link_invoices (
            id, merchant_id, payment_link_id, amount_usd, currency, network,
            amount_crypto, status, deposit_address, qr_code_data, expires_at, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING', $8, $9, $10, NOW())
        RETURNING id, merchant_id, payment_link_id, amount_usd, currency, network,
                  amount_crypto, status, deposit_address, qr_code_data, expires_at, confirmed_at, created_at
    `

    var inv mo.Invoice
    var plID sql.NullString
    var confirmedAt sql.NullTime

    err := r.db.QueryRow(
        query,
        invoiceID, merchantID, paymentLinkId, amountUSD, currency, network,
        amountCrypto, depositAddress, qrCodeData, expiresAt,
    ).Scan(
        &inv.ID, &inv.MerchantID, &plID, &inv.AmountUSD,
        &inv.Currency, &inv.Network, &inv.AmountCrypto, &inv.Status,
        &inv.DepositAddress, &inv.QRCodeData, &inv.ExpiresAt, &confirmedAt, &inv.CreatedAt,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to insert payment link invoice: %v", err)
    }

    if plID.Valid {
        inv.PaymentLinkID = &plID.String
    }
    if confirmedAt.Valid {
        inv.ConfirmedAt = &confirmedAt.Time
    }

    return &inv, nil
}

// ==========================================
// 3. CREATE DIRECT INVOICE (E-Commerce API)
// ==========================================
func (r *InvoiceRepository) CreateDirectInvoice(
    invoiceID string,
    merchantID string,
    orderID string,
    amountUSD float64,
    currency string,
    network string,
    amountCrypto float64,
    depositAddress string,
    qrCodeData string,
    expiresAt time.Time,
) (*mo.Invoice, error) {
    query := `
        INSERT INTO direct_invoices (
            id, merchant_id, order_id, amount_usd, currency, network,
            amount_crypto, status, deposit_address, qr_code_data, expires_at, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING', $8, $9, $10, NOW())
        RETURNING id, merchant_id, order_id, amount_usd, currency, network,
                  amount_crypto, status, deposit_address, qr_code_data, expires_at, confirmed_at, created_at
    `

    var inv mo.Invoice
    var orderIDCol sql.NullString
    var confirmedAt sql.NullTime

    err := r.db.QueryRow(
        query,
        invoiceID, merchantID, orderID, amountUSD, currency, network,
        amountCrypto, depositAddress, qrCodeData, expiresAt,
    ).Scan(
        &inv.ID, &inv.MerchantID, &orderIDCol, &inv.AmountUSD,
        &inv.Currency, &inv.Network, &inv.AmountCrypto, &inv.Status,
        &inv.DepositAddress, &inv.QRCodeData, &inv.ExpiresAt, &confirmedAt, &inv.CreatedAt,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to insert direct invoice: %v", err)
    }

    if confirmedAt.Valid {
        inv.ConfirmedAt = &confirmedAt.Time
    }

    return &inv, nil
}

// ==========================================
// 4. FIND INVOICE ACROSS BOTH TABLES (Unified Lookup)
// ==========================================
func (r *InvoiceRepository) FindInvoiceAnyType(invoiceID string) (*mo.Invoice, error) {
    query := `
        SELECT id, merchant_id, payment_link_id, NULL AS order_id, amount_usd, currency, network, amount_crypto, status, deposit_address, qr_code_data, expires_at, confirmed_at, created_at
        FROM payment_link_invoices
        WHERE id = $1
        UNION ALL
        SELECT id, merchant_id, NULL AS payment_link_id, order_id, amount_usd, currency, network, amount_crypto, status, deposit_address, qr_code_data, expires_at, confirmed_at, created_at
        FROM direct_invoices
        WHERE id = $1
    `

    var inv mo.Invoice
    var plID, orderID sql.NullString
    var confirmedAt sql.NullTime

    err := r.db.QueryRow(query, invoiceID).Scan(
        &inv.ID, &inv.MerchantID, &plID, &orderID, &inv.AmountUSD,
        &inv.Currency, &inv.Network, &inv.AmountCrypto, &inv.Status,
        &inv.DepositAddress, &inv.QRCodeData, &inv.ExpiresAt, &confirmedAt, &inv.CreatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("invoice not found")
        }
        return nil, err
    }

    if plID.Valid {
        inv.PaymentLinkID = &plID.String
    }
    if confirmedAt.Valid {
        inv.ConfirmedAt = &confirmedAt.Time
    }

    return &inv, nil
}
// ==========================================
// GET DIRECT INVOICE BY ID
// ==========================================
func (r *InvoiceRepository) GetDirectInvoiceByID(invoiceID string) (*mo.Invoice, error) {
	query := `
		SELECT
			id,
			merchant_id,
			order_id,
			amount_usd,
			currency,
			network,
			amount_crypto,
			status,
			deposit_address,
			qr_code_data,
			expires_at,
			confirmed_at,
			created_at
		FROM direct_invoices
		WHERE id = $1
	`

	var inv mo.Invoice
	// var orderID sql.NullString
	var confirmedAt sql.NullTime

	err := r.db.QueryRow(query, invoiceID).Scan(
		&inv.ID,
		&inv.MerchantID,
		&inv.OrderID,
		&inv.AmountUSD,
		&inv.Currency,
		&inv.Network,
		&inv.AmountCrypto,
		&inv.Status,
		&inv.DepositAddress,
		&inv.QRCodeData,
		&inv.ExpiresAt,
		&confirmedAt,
		&inv.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("direct invoice not found")
		}
		return nil, fmt.Errorf("failed to get direct invoice: %v", err)
	}



	if confirmedAt.Valid {
		inv.ConfirmedAt = &confirmedAt.Time
	}

	return &inv, nil
}
