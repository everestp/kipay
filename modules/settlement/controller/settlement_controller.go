package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-backend/modules/settlement/service"
	"go-backend/pkg/utils"
)

type SettlementController struct {
	settlementService *service.SettlementService
}

// SubmitTransactionRequest aligns 1:1 with settlement.proto payload
type VerifySettlementRequest struct {
    InvoiceID      string  `json:"invoice_id"`
	OrderID        string   `json:"order_id"`
    TxHash         string  `json:"tx_hash"`
    Network        string  `json:"network"`
    AmountPaid     float64 `json:"amount_paid"`
    Currency       string  `json:"currency"`
    SenderAddress  string  `json:"sender_address"`
    ReceiverAddress string `json:"receiver_address"`
    BlockNumber    int64   `json:"block_number"`

}
func NewSettlementController(settlementService *service.SettlementService) *SettlementController {
	return &SettlementController{
		settlementService: settlementService,
	}
}

// HandleVerificationEvent handles the transaction submission, validation, and settlement pipeline
func (c *SettlementController) HandleLinkInvoiceVerificationEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req VerifySettlementRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondValidationError(
			w,
			"Invalid JSON payload format",
		)
		return
	}

	// Basic validation
	if req.InvoiceID == "" ||
	req.OrderID == "" ||
		req.TxHash == "" ||
		req.Network == "" ||
		req.Currency == "" ||
		req.SenderAddress == "" ||
		req.ReceiverAddress == "" ||
		req.AmountPaid <= 0 {

		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Missing required transaction fields",
		)
		return
	}

	merchantID, err := c.settlementService.ProcessLinkInVoiceSettlement(
		r.Context(),
		req.InvoiceID,
		req.TxHash,
		req.Network,
		req.AmountPaid,
		req.Currency,
		req.SenderAddress,
		req.ReceiverAddress,
		req.BlockNumber,
	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrReplayDetected),
			errors.Is(err, service.ErrReplayAttack):

			utils.ErrorResponse(
				w,
				http.StatusConflict,
				"Replay attack detected: transaction has already been processed",
			)

		case errors.Is(err, service.ErrBlacklisted):

			utils.ErrorResponse(
				w,
				http.StatusForbidden,
				"Transaction rejected: address is blacklisted",
			)

		case errors.Is(err, service.ErrInvalidSignature),
			errors.Is(err, service.ErrVerificationFail):

			utils.ErrorResponse(
				w,
				http.StatusUnauthorized,
				"Blockchain transaction verification failed",
			)

		default:

			utils.ErrorResponse(
				w,
				http.StatusInternalServerError,
				err.Error(),
			)
		}

		return
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"Transaction verified and settled successfully",
		map[string]interface{}{
			"tx_hash":          req.TxHash,
			"invoice_id":       req.InvoiceID,
			"merchant_id":      merchantID,
			"sender_address":   req.SenderAddress,
			"receiver_address": req.ReceiverAddress,
			"block_number":     req.BlockNumber,
			"amount_paid":      req.AmountPaid,
			"currency":         req.Currency,
			"network":          req.Network,
			"status":            "SUCCESS",
		},
	)
}
func (c *SettlementController) HandleAPIInvpiceVerificationEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req VerifySettlementRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondValidationError(
			w,
			"Invalid JSON payload format",
		)
		return
	}

	// Basic validation
	if req.InvoiceID == "" ||
		req.TxHash == "" ||
		req.Network == "" ||
		req.Currency == "" ||
		req.SenderAddress == "" ||
		req.ReceiverAddress == "" ||
		req.AmountPaid <= 0 {

		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Missing required transaction fields",
		)
		return
	}

	merchantID, err := c.settlementService.ProcessAPIInVoiceSettlement(
		r.Context(),
		req.InvoiceID,
		req.OrderID,
		req.TxHash,
		req.Network,
		req.AmountPaid,
		req.Currency,
		req.SenderAddress,
		req.ReceiverAddress,
		req.BlockNumber,


	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrReplayDetected),
			errors.Is(err, service.ErrReplayAttack):

			utils.ErrorResponse(
				w,
				http.StatusConflict,
				"Replay attack detected: transaction has already been processed",
			)

		case errors.Is(err, service.ErrBlacklisted):

			utils.ErrorResponse(
				w,
				http.StatusForbidden,
				"Transaction rejected: address is blacklisted",
			)

		case errors.Is(err, service.ErrInvalidSignature),
			errors.Is(err, service.ErrVerificationFail):

			utils.ErrorResponse(
				w,
				http.StatusUnauthorized,
				"Blockchain transaction verification failed",
			)

		default:

			utils.ErrorResponse(
				w,
				http.StatusInternalServerError,
				err.Error(),
			)
		}

		return
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"Transaction verified and settled successfully",
		map[string]interface{}{
			"tx_hash":          req.TxHash,
			"invoice_id":       req.InvoiceID,
			"merchant_id":      merchantID,
			"sender_address":   req.SenderAddress,
			"receiver_address": req.ReceiverAddress,
			"block_number":     req.BlockNumber,
			"amount_paid":      req.AmountPaid,
			"currency":         req.Currency,
			"network":          req.Network,
			"status":            "SUCCESS",
		},
	)
}
