package controller

import (
	"fmt"
	"net/http"

	"go-backend/modules/dashboard/service"
	"go-backend/pkg/middleware"
	"go-backend/pkg/utils"
)

type DashboardController struct {
	service *service.DashboardService
}

func NewDashboardController(service *service.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

func (c *DashboardController) GetMetrics(w http.ResponseWriter, r *http.Request) {
    // FIX: Use the middleware helper function to safely fetch using the correct typed context key
    merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())

    fmt.Printf("This is in the controller: %v\n", merchantID)

    if !ok || merchantID == "" {
        utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
        return
    }

    res, err := c.service.GetMetrics(merchantID)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusOK, "Dashboard metrics retrieved successfully", res)
}
