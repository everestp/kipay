package repository

import (
	"database/sql"
	"fmt"
	"time"

	"go-backend/modules/kyc/dto"
)

type KycRepository struct {
	db *sql.DB
}

func NewKycRepository(db *sql.DB) *KycRepository {
	return &KycRepository{db: db}
}

func (r *KycRepository) UpsertDocument(docID string, merchantID string, docType string, fileUrl string) error {
	query := `
		INSERT INTO merchant_kyc_documents (id, merchant_id, doc_type, file_url, status, created_at)
		VALUES ($1, $2, $3, $4, 'PENDING', NOW())
		ON CONFLICT (id) DO UPDATE
		SET file_url = EXCLUDED.file_url, status = 'PENDING', rejection_reason = NULL
	`
	_, err := r.db.Exec(query, docID, merchantID, docType, fileUrl)
	if err != nil {
		return fmt.Errorf("failed to save kyc document: %v", err)
	}

	// Update merchant status to IN_REVIEW if currently PENDING_KYC
	_, err = r.db.Exec(`UPDATE merchants SET status = 'IN_REVIEW', kyc_submitted_at = NOW() WHERE id = $1 AND status = 'PENDING_KYC'`, merchantID)
	return err
}

func (r *KycRepository) GetMerchantKycDetails(merchantID string) (*dto.KycStatusResponse, error) {
	var mStatus string
	var subAt, verAt sql.NullTime

	err := r.db.QueryRow(`SELECT status, kyc_submitted_at, verified_at FROM merchants WHERE id = $1`, merchantID).Scan(&mStatus, &subAt, &verAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`SELECT id, doc_type, file_url, status, rejection_reason, created_at FROM merchant_kyc_documents WHERE merchant_id = $1`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []dto.KycDocItem
	for rows.Next() {
		var d dto.KycDocItem
		var rej sql.NullString
		var cAt time.Time
		if err := rows.Scan(&d.ID, &d.DocType, &d.FileUrl, &d.Status, &rej, &cAt); err == nil {
			if rej.Valid {
				d.RejectionReason = &rej.String
			}
			d.CreatedAt = cAt.Format(time.RFC3339)
			docs = append(docs, d)
		}
	}

	res := &dto.KycStatusResponse{
		MerchantID:     merchantID,
		MerchantStatus: mStatus,
		Documents:      docs,
	}
	if subAt.Valid {
		s := subAt.Time.Format(time.RFC3339)
		res.KycSubmittedAt = &s
	}
	if verAt.Valid {
		v := verAt.Time.Format(time.RFC3339)
		res.VerifiedAt = &v
	}

	return res, nil
}
