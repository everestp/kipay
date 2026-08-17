package controller

import (
	"encoding/json"
	"net/http"

	"go-backend/modules/invoice/dto"
	"go-backend/modules/invoice/service"
	"go-backend/pkg/middleware"
	"go-backend/pkg/utils"

	"github.com/gorilla/mux"
)

type InvoiceController struct {
    invoiceService *service.InvoiceService
}

func NewInvoiceController(invoiceService *service.InvoiceService) *InvoiceController {
    return &InvoiceController{invoiceService: invoiceService}
}

// ==========================================
// 1. CREATE LINK INVOICE (Public Customer Flow)
// Called when a customer opens a static Payment Link and selects a currency/network.
// No merchant session is required here because the customer is external.
// ==========================================
func (c *InvoiceController) CreateLinkInvoice(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateLinkInvoiceRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    res, err := c.invoiceService.CreateLinkInvoice(req)
    if err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusCreated, "Payment link invoice generated successfully", res)
}

// ==========================================
// 2. CREATE DIRECT INVOICE (E-Commerce API Flow)
// Called by external e-commerce plugins (Shopify, WooCommerce, custom carts).
// Authenticated via API Key middleware which injects the merchant_id into the request context.
// ==========================================
func (c *InvoiceController) CreateDirectInvoice(w http.ResponseWriter, r *http.Request) {
    // Extract merchant ID injected from API Key authentication middleware
        merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())

    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized API key or session")
        return
    }

    var req dto.CreateDirectInvoiceRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    res, err := c.invoiceService.CreateDirectInvoice(merchantID, req)
    if err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusCreated, "Direct e-commerce invoice generated successfully", res)
}
func (c *InvoiceController) CreateAPIInvoice(w http.ResponseWriter, r *http.Request) {
    // Extract merchant ID injected from API Key authentication middleware
        merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())

    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized API key or session")
        return
    }

    var req dto.CreateDirectInvoiceRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    res, err := c.invoiceService.CreateDirectInvoice(merchantID, req)
    if err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusCreated, "Direct e-commerce invoice generated successfully", res)
}

// ==========================================
// 3. GET INVOICE STATUS (Public Polling Flow)
// Called by the checkout frontend UI every few seconds to check if payment
// has been detected and confirmed on-chain by the Rust engine.
// ==========================================
func (c *InvoiceController) GetInvoiceStatus(w http.ResponseWriter, r *http.Request) {
    // Extract invoice ID from URL path parameters (e.g., chi or gorilla/mux router)
    invoiceID := r.URL.Query().Get("id") // Or chi.URLParam(r, "id") depending on your router
    if invoiceID == "" {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invoice ID is required")
        return
    }

    res, err := c.invoiceService.GetInvoiceStatus(invoiceID)
    if err != nil {
        utils.ErrorResponse(w, http.StatusNotFound, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusOK, "Invoice status fetched successfully", res)
}


// ==========================================
// GET DIRECT INVOICE BY INVOICE ID
// GET /invoices/direct/{invoiceID}
// ==========================================
func (c *InvoiceController) GetDirectInvoiceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoiceID"]

	if invoiceID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invoice ID is required")
		return
	}

	res, err := c.invoiceService.GetDirectInvoiceByID(invoiceID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"Direct invoice fetched successfully",
		res,
	)
}
