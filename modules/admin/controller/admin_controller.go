package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go-backend/modules/admin/dto"
	"go-backend/modules/admin/service"
	"go-backend/pkg/utils"
)

type AdminController struct {
	service *service.AdminService
}

func NewAdminController(service *service.AdminService) *AdminController {
	return &AdminController{service: service}
}

func (c *AdminController) ListMerchants(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	list, err := c.service.GetAllMerchants(statusFilter)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Merchants list retrieved", list)
}

func (c *AdminController) UpdateMerchantStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	merchantID := vars["merchantId"]

	var req dto.UpdateMerchantStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := c.service.ModifyMerchantStatus(merchantID, req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Merchant status updated and notification dispatched successfully", nil)
}
