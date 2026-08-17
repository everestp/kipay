package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"go-backend/modules/wallet/dto"
	"go-backend/modules/wallet/service"
	"go-backend/pkg/middleware"
	"go-backend/pkg/utils"
)

type WalletController struct {
	service *service.WalletService
}

func NewWalletController(service *service.WalletService) *WalletController {
	return &WalletController{service: service}
}

// Create Wallet
func (c *WalletController) Create(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	var req dto.CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	res, err := c.service.CreateWallet(merchantID, req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusCreated,
		"Wallet created successfully",
		res,
	)
}

// List All Wallets for Merchant
func (c *WalletController) List(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	wallets, err := c.service.ListMerchantWallets(merchantID)
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
		"Wallets retrieved successfully",
		wallets,
	)
}

// Get Wallet by ID
func (c *WalletController) GetByID(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	walletID, err := getWalletID(r)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid wallet ID",
		)
		return
	}

	wallet, err := c.service.GetWalletByID(merchantID, walletID)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"Wallet retrieved successfully",
		wallet,
	)
}

// Update Wallet
func (c *WalletController) Update(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	walletID, err := getWalletID(r)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid wallet ID",
		)
		return
	}

	var req dto.UpdateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
		)
		return
	}

	res, err := c.service.UpdateWallet(merchantID, walletID, req)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"Wallet updated successfully",
		res,
	)
}

// Delete Wallet
func (c *WalletController) Delete(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized session")
		return
	}

	walletID, err := getWalletID(r)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid wallet ID",
		)
		return
	}

	err = c.service.DeleteWallet(merchantID, walletID)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		"Wallet deleted successfully",
		nil,
	)
}

// getWalletID parses the wallet ID from the URL path.
func getWalletID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)

	id := vars["id"]
	if id == "" {
		return 0, fmt.Errorf("wallet ID path parameter is required")
	}

	return strconv.ParseInt(id, 10, 64)
}
