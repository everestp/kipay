package controller

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "go-backend/modules/payment_link/dto"
    "go-backend/modules/payment_link/service"
    "go-backend/pkg/middleware"
    "go-backend/pkg/utils"
)

type PaymentLinkController struct {
    service *service.PaymentLinkService
}

func NewPaymentLinkController(service *service.PaymentLinkService) *PaymentLinkController {
    return &PaymentLinkController{service: service}
}

func (c *PaymentLinkController) Create(w http.ResponseWriter, r *http.Request) {
    merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
        return
    }

    var req dto.CreatePaymentLinkRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    res, err := c.service.CreatePaymentLink(merchantID, req)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusCreated, "Payment link created successfully", res)
}

func (c *PaymentLinkController) GetPAyentLinkById(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    linkID := vars["linkId"]

    res, err := c.service.GetPaymentLink(linkID)
    if err != nil {
        utils.ErrorResponse(w, http.StatusNotFound, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusOK, "Payment link retrieved", res)
}

func (c *PaymentLinkController) GetAllPaymentLinkByMerchant(w http.ResponseWriter, r *http.Request) {
    merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
        return
    }

    links, err := c.service.GetAllPaymentLinksByMerchant(merchantID)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusOK, "Payment links retrieved successfully", links)
}

// Update handles modifying an existing payment link
func (c *PaymentLinkController) Update(w http.ResponseWriter, r *http.Request) {
    merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
        return
    }

    vars := mux.Vars(r)
    paymentLinkID := vars["linkId"] // Matches route parameter like /payment-links/{linkId}

    var req dto.UpdatePaymentLinkRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    res, err := c.service.UpdatePaymentLink(paymentLinkID, merchantID, req)
    if err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusOK, "Payment link updated successfully", res)
}
