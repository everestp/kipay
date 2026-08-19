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
	ErrReplayAttack     = ErrReplayDetected
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

// ============================================================
// LINK INVOICE SETTLEMENT
// ============================================================

func (s *SettlementService) ProcessLinkInVoiceSettlement(
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
		s.rustClient.VerifyAndSettleLinkInvoiceTransaction(
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
			"[Settlement] Link invoice verification failed | invoice=%s | tx=%s | message=%s",
			invoiceID,
			txHash,
			message,
		)

		return "", ErrVerificationFail
	}

	// =========================================================
	// 4. RECORD VERIFIED SETTLEMENT
	// =========================================================

	merchantID, err := s.repo.RecordLinkInvoiceSettlement(
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

	// =========================================================
	// 6. DISPATCH WEBHOOK
	// =========================================================

	s.dispatchInvoiceConfirmedWebhook(
		merchantID,
		invoiceID,
		txHash,
		network,
		amountPaid,
		currency,
		senderAddress,
		receiverAddress,
		blockNumber,
	)

	// =========================================================
	// 7. SUCCESS LOG
	// =========================================================

	log.Printf(
		"[Settlement] LINK SUCCESS | tx=%s | invoice=%s | merchant=%s | sender=%s | receiver=%s | amount=%f %s | block=%d",
		txHash,
		invoiceID,
		merchantID,
		senderAddress,
		receiverAddress,
		amountPaid,
		currency,
		blockNumber,
	)

	return merchantID, nil
}

// ============================================================
// API INVOICE SETTLEMENT
// ============================================================

func (s *SettlementService) ProcessAPIInVoiceSettlement(
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
		s.rustClient.VerifyAndSettleAPIInvoiceTransaction(
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
			"[Settlement] API invoice verification failed | invoice=%s | tx=%s | message=%s",
			invoiceID,
			txHash,
			message,
		)

		return "", ErrVerificationFail
	}

	// =========================================================
	// 4. RECORD VERIFIED SETTLEMENT
	// =========================================================

	merchantID, err := s.repo.RecordAPIInvoiceSettlement(
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

	// =========================================================
	// 6. DISPATCH WEBHOOK
	// =========================================================

	s.dispatchInvoiceConfirmedWebhook(
		merchantID,
		invoiceID,
		txHash,
		network,
		amountPaid,
		currency,
		senderAddress,
		receiverAddress,
		blockNumber,
	)

	// =========================================================
	// 7. SUCCESS LOG
	// =========================================================

	log.Printf(
		"[Settlement] API SUCCESS | tx=%s | invoice=%s | merchant=%s | sender=%s | receiver=%s | amount=%f %s | block=%d",
		txHash,
		invoiceID,
		merchantID,
		senderAddress,
		receiverAddress,
		amountPaid,
		currency,
		blockNumber,
	)

	return merchantID, nil
}

// ============================================================
// INVOICE CONFIRMED WEBHOOK
// ============================================================

func (s *SettlementService) dispatchInvoiceConfirmedWebhook(
	merchantID string,
	invoiceID string,
	txHash string,
	network string,
	amountPaid float64,
	currency string,
	senderAddress string,
	receiverAddress string,
	blockNumber int64,
) {

	/*
	 * Payment settlement has already been committed
	 * successfully at this point.
	 *
	 * Therefore webhook failure MUST NOT cause the
	 * payment settlement itself to fail.
	 */

	if s.webhookSvc == nil {
		log.Printf(
			"[Webhook] Webhook service is nil | merchant=%s | invoice=%s",
			merchantID,
			invoiceID,
		)

		return
	}

	payload := map[string]interface{}{
		"invoice_id":     invoiceID,
		"transaction_id": txHash,
		"tx_hash":       txHash,
		"merchant_id":   merchantID,
		"network":       network,
		"amount_paid":   amountPaid,
		"currency":      currency,
		"sender_address": senderAddress,
		"receiver_address": receiverAddress,
		"block_number":  blockNumber,
		"status":        "CONFIRMED",
	}

	err := s.webhookSvc.DispatchEvent(
		merchantID,
		"invoice.confirmed",
		payload,
	)

	if err != nil {
		log.Printf(
			"[Webhook] Failed to dispatch invoice.confirmed | merchant=%s | invoice=%s | tx=%s | error=%v",
			merchantID,
			invoiceID,
			txHash,
			err,
		)

		return
	}

	log.Printf(
		"[Webhook] invoice.confirmed dispatched successfully | merchant=%s | invoice=%s | tx=%s",
		merchantID,
		invoiceID,
		txHash,
	)
}
