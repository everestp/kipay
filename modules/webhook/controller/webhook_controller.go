package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"go-backend/modules/webhook/dto"
	"go-backend/modules/webhook/service"
	"go-backend/pkg/middleware"
	"go-backend/pkg/utils"

	"github.com/gorilla/mux"
)

type WebhookController struct {
	service *service.WebhookService
}

func NewWebhookController(
	service *service.WebhookService,
) *WebhookController {
	return &WebhookController{
		service: service,
	}
}

// ============================================================
// REGISTER WEBHOOK
// POST /api/v1/webhooks
// ============================================================

func (c *WebhookController) Register(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get authenticated merchant ID
	// ----------------------------------------------------------

 merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())

	if !ok || strings.TrimSpace(merchantID) == "" {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"Unauthorized session",
		)
		return
	}

	// ----------------------------------------------------------
	// Decode request
	// ----------------------------------------------------------

	var req dto.RegisterWebhookEndpointRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
		)
		return
	}

	// ----------------------------------------------------------
	// Validate endpoint URL
	// ----------------------------------------------------------

	req.EndpointURL =
		strings.TrimSpace(req.EndpointURL)

	if req.EndpointURL == "" {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Endpoint URL is required",
		)
		return
	}

	// ----------------------------------------------------------
	// Validate subscribed events
	// ----------------------------------------------------------

	if len(req.SubscribedEvents) == 0 {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"At least one event must be subscribed",
		)
		return
	}

	// ----------------------------------------------------------
	// Register webhook
	//
	// Service returns:
	//
	// endpointID
	// secret
	// error
	// ----------------------------------------------------------

	endpointID, secret, err :=
		c.service.RegisterEndpoint(
			merchantID,
			req,
		)

	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	// ----------------------------------------------------------
	// Return secret ONLY during creation
	// ----------------------------------------------------------

	utils.SuccessResponse(
		w,
		http.StatusCreated,
		"Webhook endpoint registered successfully",
		map[string]interface{}{
			"id":                 endpointID,
			"endpoint_url":       req.EndpointURL,
			"subscribed_events":  req.SubscribedEvents,
			"secret":             secret,
			"active":             true,
		},
	)
}

// ============================================================
// LIST WEBHOOKS
// GET /api/v1/webhooks
// ============================================================

func (c *WebhookController) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get merchant ID
	// ----------------------------------------------------------

 merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())

	if !ok || strings.TrimSpace(merchantID) == "" {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"Unauthorized session",
		)
		return
	}

	// ----------------------------------------------------------
	// Get endpoints
	// ----------------------------------------------------------

	webhooks, err :=
		c.service.ListEndpoints(
			merchantID,
		)

	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to fetch webhook endpoints",
		)
		return
	}

	// ----------------------------------------------------------
	// Return endpoints
	//
	// IMPORTANT:
	// List should NOT return webhook secrets.
	// ----------------------------------------------------------

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"Webhook endpoints fetched successfully",
		webhooks,
	)
}

// ============================================================
// DELETE / DEACTIVATE WEBHOOK
// DELETE /api/v1/webhooks/{id}
// ============================================================

func (c *WebhookController) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get merchant ID
	// ----------------------------------------------------------

 merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())

	if !ok || strings.TrimSpace(merchantID) == "" {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"Unauthorized session",
		)
		return
	}

	// ----------------------------------------------------------
	// Get endpoint ID
	// ----------------------------------------------------------

	vars := mux.Vars(r)

	endpointID :=
		strings.TrimSpace(
			vars["id"],
		)

	if endpointID == "" {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Webhook endpoint ID is required",
		)
		return
	}

	// ----------------------------------------------------------
	// Deactivate endpoint
	// ----------------------------------------------------------

	err :=
		c.service.DeactivateEndpoint(
			merchantID,
			endpointID,
		)

	if err != nil {

		if errors.Is(
			err,
			service.ErrWebhookNotFound,
		) {
			utils.ErrorResponse(
				w,
				http.StatusNotFound,
				"Webhook endpoint not found",
			)
			return
		}

		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to delete webhook endpoint",
		)
		return
	}

	// ----------------------------------------------------------
	// Success
	// ----------------------------------------------------------

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"Webhook endpoint deleted successfully",
		map[string]interface{}{
			"id":     endpointID,
			"active": false,
		},
	)
}

// ============================================================
// ROTATE SECRET
// POST /api/v1/webhooks/{id}/rotate-secret
// ============================================================

func (c *WebhookController) RotateSecret(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get merchant ID
	// ----------------------------------------------------------

 merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())

	if !ok || strings.TrimSpace(merchantID) == "" {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"Unauthorized session",
		)
		return
	}

	// ----------------------------------------------------------
	// Get endpoint ID
	// ----------------------------------------------------------

	vars := mux.Vars(r)

	endpointID :=
		strings.TrimSpace(
			vars["id"],
		)

	if endpointID == "" {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Webhook endpoint ID is required",
		)
		return
	}

	// ----------------------------------------------------------
	// Rotate secret
	// ----------------------------------------------------------

	secret, err :=
		c.service.RotateSecret(
			merchantID,
			endpointID,
		)

	if err != nil {

		if errors.Is(
			err,
			service.ErrWebhookNotFound,
		) {
			utils.ErrorResponse(
				w,
				http.StatusNotFound,
				"Webhook endpoint not found",
			)
			return
		}

		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to rotate webhook secret",
		)
		return
	}

	// ----------------------------------------------------------
	// Return new secret
	//
	// The merchant should save this immediately.
	// ----------------------------------------------------------

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"Webhook secret rotated successfully",
		map[string]interface{}{
			"id":     endpointID,
			"secret": secret,
		},
	)
}
type TestWebhookRequest struct {
	ID string `json:"id"`
}

type TestWebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c *WebhookController) TestWebhook(
	w http.ResponseWriter,
	r *http.Request,
) {
        merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())




	if !ok || merchantID == "" {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

vars := mux.Vars(r)

webhookID := vars["id"]

if webhookID == "" {
	http.Error(
		w,
		"webhook id is required",
		http.StatusBadRequest,
	)
	return
}

	if webhookID == "" {
		http.Error(
			w,
			"webhook id is required",
			http.StatusBadRequest,
		)
		return
	}

	err :=
		c.service.SendTestWebhook(
			r.Context(),
			webhookID,
			merchantID,
		)

	if err != nil {
		response := TestWebhookResponse{
			Success: false,
			Message: err.Error(),
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		w.WriteHeader(
			http.StatusBadGateway,
		)

		_ = json.NewEncoder(w).Encode(
			response,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(
		TestWebhookResponse{
			Success: true,
			Message: "Test webhook sent successfully",
		},
	)
}
