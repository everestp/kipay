package controller

import (
	"encoding/json"
	"net/http"

	"go-backend/modules/webhook/dto"
	"go-backend/modules/webhook/service"
	"go-backend/pkg/utils"
)

type WebhookController struct {
	service *service.WebhookService
}

func NewWebhookController(service *service.WebhookService) *WebhookController {
	return &WebhookController{service: service}
}

func (c *WebhookController) Register(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := r.Context().Value("merchant_id").(string)
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	var req dto.RegisterWebhookEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := c.service.RegisterEndpoint(merchantID, req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Webhook endpoint registered successfully", nil)
}
