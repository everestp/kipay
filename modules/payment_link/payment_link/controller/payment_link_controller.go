package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go-backend/modules/payment_link/dto"
	"go-backend/modules/payment_link/service"
	"go-backend/pkg/utils"
)

type PaymentLinkController struct {
	service *service.PaymentLinkService
}

func NewPaymentLinkController(service *service.PaymentLinkService) *PaymentLinkController {
	return &PaymentLinkController{service: service}
}

func (c *PaymentLinkController) Create(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := r.Context().Value("merchant_id").(string)
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

func (c *PaymentLinkController) GetByPublicID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkID := vars["linkId"]

	res, err := c.service.GetPaymentLink(linkID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Payment link retrieved", res)
}
