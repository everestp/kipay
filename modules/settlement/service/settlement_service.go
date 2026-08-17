package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go-backend/modules/settlement/client"
	"go-backend/modules/settlement/repository"
	webhookService "go-backend/modules/webhook/service"
)

var (
	ErrReplayDetected   = errors.New("replay attack detected: transaction hash has already been processed")
	ErrReplayAttack     = ErrReplayDetected // Alias for backward compatibility
	ErrBlacklisted      = errors.New("sender or recipient address is blacklisted")
	ErrInvalidSignature = errors.New("invalid transaction signature verified by security engine")
	ErrVerificationFail = ErrInvalidSignature
)

type SettlementService struct {
	repo       *repository.SettlementRepository
	webhookSvc *webhookService.WebhookService
	rustClient *client.RustVerificationClient
}

func NewSettlementService(
	repo *repository.SettlementRepository,
	webhookSvc *webhookService.WebhookService,
	rustClient *client.RustVerificationClient,
) *SettlementService {
	return &SettlementService{
		repo:       repo,
		webhookSvc: webhookSvc,
		rustClient: rustClient,
	}
}
// ProcessSettlement orchestrates security pre-checks, Rust gRPC crypto verification, database persistence, and webhooks
func (s *SettlementService) ProcessSettlement(
	ctx context.Context,
	invoiceID string,
	txHash string,
	network string,
	amountPaid float64,
	currency string,
	senderAddress string,
	receiverAddress string,
	blockNumber int64,
) (string, error) {
	// =========================================================
	// 1. BLACKLIST CHECK
	// =========================================================

	senderBlocked, err :=
		s.repo.IsAddressBlacklisted(
			ctx,
			senderAddress,
		)

	if err != nil {
		return "",
			fmt.Errorf(
				"failed to verify sender blacklist status: %w",
				err,
			)
	}

	if senderBlocked {
		return "", ErrBlacklisted
	}

	// =========================================================
	// 2. REPLAY PROTECTION
	// =========================================================

	exists, err :=
		s.repo.IsTxProcessed(
			ctx,
			txHash,
		)

	if err != nil {
		return "",
			fmt.Errorf(
				"failed to check transaction replay status: %w",
				err,
			)
	}

	if exists {
		return "", ErrReplayDetected
	}

	// =========================================================
	// 3. VERIFY WITH RUST
	// =========================================================

	success,
		resolvedMerchantID,
		message,
		err :=
		s.rustClient.VerifyAndSettleTransaction(
			ctx,
			invoiceID,
			txHash,
			network,
			amountPaid,
			currency,
			senderAddress,
			receiverAddress,
			blockNumber,
		)

	if err != nil {
		return "",
			fmt.Errorf(
				"rust verification service unavailable: %w",
				err,
			)
	}

	if !success {
		log.Printf(
			"[Settlement] Verification failed from Rust engine: %s",
			message,
		)

		return "", ErrVerificationFail
	}

	// =========================================================
	// 4. RECORD VERIFIED SETTLEMENT
	// =========================================================

	merchantID, err := s.repo.RecordSettlement(
	ctx,
	invoiceID,
	txHash,
	senderAddress,
	receiverAddress,
	amountPaid,
	currency,
	network,
	blockNumber,
)

	if err != nil {
		return "",
			fmt.Errorf(
				"failed to record settlement: %w",
				err,
			)
	}

	// =========================================================
	// 5. USE MERCHANT ID FROM RUST IF AVAILABLE
	// =========================================================

	if resolvedMerchantID != "" {
		merchantID = resolvedMerchantID
	}

	log.Printf(
		"[Settlement] SUCCESS | tx=%s | invoice=%s | sender=%s | receiver=%s | amount=%f %s | block=%d",
		txHash,
		invoiceID,
		senderAddress,
		receiverAddress,
		amountPaid,
		currency,
		blockNumber,
	)

	return merchantID, nil
}
