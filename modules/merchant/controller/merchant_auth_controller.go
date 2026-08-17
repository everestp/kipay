package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"go-backend/modules/merchant/dto"
	"go-backend/modules/merchant/service"
	"go-backend/pkg/middleware"
	"go-backend/pkg/utils"
)

type MerchantAuthController struct {
    authService *service.MerchantAuthService
}

func NewMerchantAuthController(authService *service.MerchantAuthService) *MerchantAuthController {
    return &MerchantAuthController{authService: authService}
}

func (c *MerchantAuthController) Register(w http.ResponseWriter, r *http.Request) {
    var req dto.RegisterMerchantRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    res, err := c.authService.Register(req)
    if err != nil {
        utils.ErrorResponse(w, http.StatusConflict, err.Error())
        return
    }

    utils.SuccessResponse(w, http.StatusCreated, "Merchant registered successfully", res)
}

func (c *MerchantAuthController) Login(w http.ResponseWriter, r *http.Request) {
    var req dto.LoginMerchantRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    res, err := c.authService.Login(req)
    if err != nil {
        utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
        return
    }

    // 1. Set the HttpOnly Cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    res.Token, // Assuming res.Token holds the JWT string
        Path:     "/",
        HttpOnly: true,                  // Blocks JS (XSS security)
        Secure:   false,                 // Set to true in Production (HTTPS)
        SameSite: http.SameSiteLaxMode,  // Protection against CSRF
        MaxAge:   3600 * 24,             // 1 day in seconds
    })

    // 2. (Optional) Strip token from the JSON body so JS never touches it directly
    // res.Token = ""

    utils.SuccessResponse(w, http.StatusOK, "Login successful", res)
}

func (c *MerchantAuthController) Logout(w http.ResponseWriter, r *http.Request) {
    // Clear the cookie by expiring it immediately
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // Set to true in Production
        SameSite: http.SameSiteLaxMode,
        MaxAge:   -1,
        Expires:  time.Unix(0, 0),
    })

    utils.SuccessResponse(w, http.StatusOK, "Logged out successfully", nil)
}
func (c *MerchantAuthController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// 1. Extract merchant_id from request context set by your authentication middleware
	merchantID, ok := middleware.GetMerchantIDFromContext(r.Context())
	if !ok || merchantID == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: merchant session missing or invalid")
		return
	}

	// 2. Call the service layer with the plain merchant ID
	res, err := c.authService.GetMeService(merchantID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	// 3. Return the sanitized profile response
	utils.SuccessResponse(w, http.StatusOK, "Current user fetched successfully", res)
}
