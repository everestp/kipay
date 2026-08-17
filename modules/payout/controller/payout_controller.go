package controller

import (
	"encoding/json"
	"net/http"

	"go-backend/modules/payout/dto"
	"go-backend/modules/payout/service"
	"go-backend/pkg/utils"
)

type PayoutController struct {
	service *service.PayoutService
}

func NewPayoutController(service *service.PayoutService) *PayoutController {
	return &PayoutController{service: service}
}

func (c *PayoutController) Create(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := r.Context().Value("merchant_id").(string)
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	var req dto.CreatePayoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	res, err := c.service.RequestPayout(merchantID, req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Payout withdrawal requested successfully", res)
}
