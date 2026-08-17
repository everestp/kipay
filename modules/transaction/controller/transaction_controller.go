package controller

import (
    "net/http"
    "strconv"

    "github.com/gorilla/mux" // 1. Import gorilla/mux for mux.Vars
    "go-backend/modules/transaction/service"
    "go-backend/pkg/middleware"
    "go-backend/pkg/utils"
)

type TransactionController struct {
    service *service.TransactionService
}

func NewTransactionController(service *service.TransactionService) *TransactionController {
    return &TransactionController{service: service}
}

func (c *TransactionController) ListInvoices(w http.ResponseWriter, r *http.Request) {
    merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
        return
    }

    statusFilter := r.URL.Query().Get("status")
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit <= 0 {
        limit = 10
    }
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

    list, err := c.service.GetMerchantInvoices(merchantID, statusFilter, limit, offset)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusOK, "Invoices list retrieved successfully", list)
}

func (c *TransactionController) ListAllLinkInvoices(w http.ResponseWriter, r *http.Request) {
    merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
        return
    }

    statusFilter := r.URL.Query().Get("status")
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit <= 0 {
        limit = 10
    }
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

    list, err := c.service.GetMerchantLinkInvoices(merchantID, statusFilter, limit, offset)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusOK, "Invoices list retrieved successfully", list)
}

func (c *TransactionController) ListAllLinkInvoicesByPaymentLinkId(w http.ResponseWriter, r *http.Request) {
    merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
        return
    }

    // 2. Extract path parameter using Gorilla Mux vars
    vars := mux.Vars(r)
    paymentLinkId := vars["payment-link-id"]

    if paymentLinkId == "" {
        utils.ErrorResponse(w, http.StatusBadRequest, "payment-link-id path parameter is required")
        return
    }

    statusFilter := r.URL.Query().Get("status")
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit <= 0 {
        limit = 10
    }
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

    list, err := c.service.GetMerchantLinkInvoicesByPaymentLinkId(merchantID, paymentLinkId, statusFilter, limit, offset)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusOK, "Invoices list retrieved successfully", list)
}

func (c *TransactionController) ListTransactions(w http.ResponseWriter, r *http.Request) {
    merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
        return
    }

    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit <= 0 {
        limit = 10
    }
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

    list, err := c.service.GetMerchantSettledTransactions(merchantID, limit, offset)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusOK, "Settled transactions list retrieved successfully", list)
}
