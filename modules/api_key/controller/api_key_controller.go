package controller

import (
	"encoding/json"
	"net/http"

	"go-backend/modules/api_key/dto"
	"go-backend/modules/api_key/service"
	"go-backend/pkg/middleware"
	"go-backend/pkg/utils"

	"github.com/gorilla/mux"
)

type ApiKeyController struct {
	service *service.ApiKeyService
}

func NewApiKeyController(service *service.ApiKeyService) *ApiKeyController {
	return &ApiKeyController{service: service}
}
func (c *ApiKeyController) Create(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	var req dto.CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	res, err := c.service.GenerateApiKey(merchantID, req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusCreated,
		"API Key generated successfully. Save your secret key now as it will not be shown again.",
		res,
	)
}
func (c *ApiKeyController) Delete(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	vars := mux.Vars(r)
	keyID := vars["id"]

	if keyID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "API key ID is required")
		return
	}

	err := c.service.DeleteApiKey(merchantID, keyID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"API Key deleted successfully",
		nil,
	)
}
func (c *ApiKeyController) GetAll(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	keys, err := c.service.GetApiKeys(merchantID)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"API keys fetched successfully",
		keys,
	)
}
