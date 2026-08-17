package repository

import (
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    "go-backend/modules/payment_link/dto"
)

type PaymentLinkRepository struct {
    db *sql.DB
}

func NewPaymentLinkRepository(db *sql.DB) *PaymentLinkRepository {
    return &PaymentLinkRepository{db: db}
}

func (r *PaymentLinkRepository) Create(linkID string, merchantID string, req dto.CreatePaymentLinkRequest) (*dto.PaymentLinkResponse, error) {
    query := `
        INSERT INTO payment_links (
            id, merchant_id, title, description, image_url, pricing_type,
            amount_usd, min_amount_usd, supported_currencies, supported_networks,
            success_message, redirect_url, continue_button_text, total_revenue_usd, is_active, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, 0.00, TRUE, NOW())
        RETURNING id, merchant_id, title, description, image_url, pricing_type,
                  amount_usd, min_amount_usd, success_message, redirect_url, continue_button_text, is_active, created_at
    `

    var res dto.PaymentLinkResponse
    var createdAtTime time.Time

    err := r.db.QueryRow(
        query,
        linkID, merchantID, req.Title, req.Description, req.ImageUrl, req.PricingType,
        req.AmountUSD, req.MinAmountUSD, pqArray(req.SupportedCurrencies), pqArray(req.SupportedNetworks),
        req.SuccessMessage, req.RedirectUrl, req.ContinueButtonText,
    ).Scan(
        &res.ID, &res.MerchantID, &res.Title, &res.Description, &res.ImageUrl, &res.PricingType,
        &res.AmountUSD, &res.MinAmountUSD, &res.SuccessMessage, &res.RedirectUrl, &res.ContinueButtonText,
        &res.IsActive, &createdAtTime,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to create payment link: %v", err)
    }

    res.SupportedCurrencies = req.SupportedCurrencies
    res.SupportedNetworks = req.SupportedNetworks
    res.CreatedAt = createdAtTime.Format(time.RFC3339)

    return &res, nil
}

func (r *PaymentLinkRepository) GetByID(id string) (*dto.PaymentLinkResponse, error) {
    query := `
       SELECT
			id,
			merchant_id,
			title,
			description,
			image_url,
			pricing_type,
			amount_usd,
			min_amount_usd,
			supported_currencies,
			supported_networks,
			success_message,
			redirect_url,
			continue_button_text,
			total_revenue_usd,
			color,
			is_active,
			created_at
		FROM payment_links
		 WHERE id = $1 AND is_active = TRUE
    `
    var res dto.PaymentLinkResponse
    var currStr, netStr string
    var createdAtTime time.Time

    err := r.db.QueryRow(query, id).Scan(
      	&res.ID,
			&res.MerchantID,
			&res.Title,
			&res.Description,
			&res.ImageUrl,
			&res.PricingType,
			&res.AmountUSD,
			&res.MinAmountUSD,
			&currStr,
			&netStr,
			&res.SuccessMessage,
			&res.RedirectUrl,
			&res.ContinueButtonText,
			&res.TotalRevenueUSD,
			&res.Color,
			&res.IsActive,
			&createdAtTime,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New("payment link not found or inactive")
    } else if err != nil {
        return nil, err
    }

    res.SupportedCurrencies = parsePqArray(currStr)
    res.SupportedNetworks = parsePqArray(netStr)
    res.CreatedAt = createdAtTime.Format(time.RFC3339)

    return &res, nil
}

func (r *PaymentLinkRepository) GetAllByMerchant(merchantID string) ([]dto.PaymentLinkResponse, error) {
	query := `
		SELECT
			id,
			merchant_id,
			title,
			description,
			image_url,
			pricing_type,
			amount_usd,
			min_amount_usd,
			supported_currencies,
			supported_networks,
			success_message,
			redirect_url,
			continue_button_text,
			total_revenue_usd,
			color,
			is_active,
			created_at
		FROM payment_links
		WHERE merchant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []dto.PaymentLinkResponse

	for rows.Next() {
		var res dto.PaymentLinkResponse
		var currStr, netStr string
		var createdAtTime time.Time

		if err := rows.Scan(
			&res.ID,
			&res.MerchantID,
			&res.Title,
			&res.Description,
			&res.ImageUrl,
			&res.PricingType,
			&res.AmountUSD,
			&res.MinAmountUSD,
			&currStr,
			&netStr,
			&res.SuccessMessage,
			&res.RedirectUrl,
			&res.ContinueButtonText,
			&res.TotalRevenueUSD,
			&res.Color,
			&res.IsActive,
			&createdAtTime,
		); err != nil {
			continue
		}

		res.SupportedCurrencies = parsePqArray(currStr)
		res.SupportedNetworks = parsePqArray(netStr)
		res.CreatedAt = createdAtTime.Format(time.RFC3339)

		list = append(list, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if list == nil {
		list = []dto.PaymentLinkResponse{}
	}

	return list, nil
}
func (r *PaymentLinkRepository) Update(paymentLinkID string, merchantID string, req dto.UpdatePaymentLinkRequest) (*dto.PaymentLinkResponse, error) {
    // Collect fields to update dynamically
    var setClauses []string
    var args []interface{}
    argIdx := 1

    // Helper function to append fields if they are provided
    if req.Title != nil {
        setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIdx))
        args = append(args, *req.Title)
        argIdx++
    }
    if req.Description != nil {
        setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
        args = append(args, *req.Description)
        argIdx++
    }
    if req.ImageUrl != nil {
        setClauses = append(setClauses, fmt.Sprintf("image_url = $%d", argIdx))
        args = append(args, *req.ImageUrl)
        argIdx++
    }
    if req.PricingType != nil {
        setClauses = append(setClauses, fmt.Sprintf("pricing_type = $%d", argIdx))
        args = append(args, *req.PricingType)
        argIdx++
    }
    if req.AmountUSD != nil {
        setClauses = append(setClauses, fmt.Sprintf("amount_usd = $%d", argIdx))
        args = append(args, *req.AmountUSD)
        argIdx++
    }
    if req.MinAmountUSD != nil {
        setClauses = append(setClauses, fmt.Sprintf("min_amount_usd = $%d", argIdx))
        args = append(args, *req.MinAmountUSD)
        argIdx++
    }
    if req.SupportedCurrencies != nil {
        setClauses = append(setClauses, fmt.Sprintf("supported_currencies = $%d", argIdx))
        args = append(args, pqArray(req.SupportedCurrencies))
        argIdx++
    }
    if req.SupportedNetworks != nil {
        setClauses = append(setClauses, fmt.Sprintf("supported_networks = $%d", argIdx))
        args = append(args, pqArray(req.SupportedNetworks))
        argIdx++
    }
    if req.SuccessMessage != nil {
        setClauses = append(setClauses, fmt.Sprintf("success_message = $%d", argIdx))
        args = append(args, *req.SuccessMessage)
        argIdx++
    }
    if req.RedirectUrl != nil {
        setClauses = append(setClauses, fmt.Sprintf("redirect_url = $%d", argIdx))
        args = append(args, *req.RedirectUrl)
        argIdx++
    }
    if req.ContinueButtonText != nil {
        setClauses = append(setClauses, fmt.Sprintf("continue_button_text = $%d", argIdx))
        args = append(args, *req.ContinueButtonText)
        argIdx++
    }
    if req.Color != nil {
        setClauses = append(setClauses, fmt.Sprintf("color = $%d", argIdx))
        args = append(args, *req.Color)
        argIdx++
    }
    if req.IsActive != nil {
        setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIdx))
        args = append(args, *req.IsActive)
        argIdx++
    }

    if len(setClauses) == 0 {
        return nil, errors.New("no fields provided for update")
    }

    // Append WHERE clause arguments for ID and MerchantID ownership validation
    query := fmt.Sprintf(`
        UPDATE payment_links
        SET %s
        WHERE id = $%d AND merchant_id = $%d
        RETURNING id, merchant_id, title, description, image_url, pricing_type,
                  amount_usd, min_amount_usd, supported_currencies, supported_networks,
                  success_message, redirect_url, continue_button_text, total_revenue_usd,
                  color, is_active, created_at
    `, strings.Join(setClauses, ", "), argIdx, argIdx+1)

    args = append(args, paymentLinkID, merchantID)

    var res dto.PaymentLinkResponse
    var currStr, netStr string
    var createdAtTime time.Time

    err := r.db.QueryRow(query, args...).Scan(
        &res.ID,
        &res.MerchantID,
        &res.Title,
        &res.Description,
        &res.ImageUrl,
        &res.PricingType,
        &res.AmountUSD,
        &res.MinAmountUSD,
        &currStr,
        &netStr,
        &res.SuccessMessage,
        &res.RedirectUrl,
        &res.ContinueButtonText,
        &res.TotalRevenueUSD,
        &res.Color,
        &res.IsActive,
        &createdAtTime,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New("payment link not found or unauthorized")
    } else if err != nil {
        return nil, fmt.Errorf("failed to update payment link: %v", err)
    }

    res.SupportedCurrencies = parsePqArray(currStr)
    res.SupportedNetworks = parsePqArray(netStr)
    res.CreatedAt = createdAtTime.Format(time.RFC3339)

    return &res, nil
}


// Helper to format string slice into PostgreSQL array syntax string format like "{USDT,USDC}"
func pqArray(items []string) string {
    return "{" + strings.Join(items, ",") + "}"
}

// Helper to parse PostgreSQL array string format back into string slice
func parsePqArray(pq string) []string {
    pq = strings.Trim(pq, "{}")
    if pq == "" {
        return []string{}
    }
    return strings.Split(pq, ",")
}
