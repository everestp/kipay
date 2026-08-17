package controller

import (
	"encoding/json"
	"net/http"

	"go-backend/modules/kyc/dto"
	"go-backend/modules/kyc/service"
	"go-backend/pkg/utils"
)

type KycController struct {
	service *service.KycService
}

func NewKycController(service *service.KycService) *KycController {
	return &KycController{service: service}
}

func (c *KycController) Submit(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := r.Context().Value("merchant_id").(string)
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	var req dto.SubmitKycDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	err := c.service.SubmitDocument(merchantID, req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "KYC document submitted successfully. Account status is now IN_REVIEW.", nil)
}

func (c *KycController) GetStatus(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := r.Context().Value("merchant_id").(string)
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	res, err := c.service.GetStatus(merchantID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "KYC status retrieved", res)
}
