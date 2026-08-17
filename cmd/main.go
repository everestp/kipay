package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"go-backend/database"
	"go-backend/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	// Admin
	adminController "go-backend/modules/admin/controller"
	adminRepo "go-backend/modules/admin/repository"
	adminService "go-backend/modules/admin/service"

	// API Key
	apiKeyController "go-backend/modules/api_key/controller"
	apiKeyRepo "go-backend/modules/api_key/repository"
	apiKeyService "go-backend/modules/api_key/service"

	// Dashboard
	dashboardController "go-backend/modules/dashboard/controller"
	dashboardRepo "go-backend/modules/dashboard/repository"
	dashboardService "go-backend/modules/dashboard/service"

	// Invoice
	invoiceController "go-backend/modules/invoice/controller"
	invoiceRepo "go-backend/modules/invoice/repository"
	invoiceService "go-backend/modules/invoice/service"

	// KYC
	kycController "go-backend/modules/kyc/controller"
	kycRepo "go-backend/modules/kyc/repository"
	kycService "go-backend/modules/kyc/service"

	// Merchant Auth
	merchantAuthController "go-backend/modules/merchant/controller"
	merchantAuthRepo "go-backend/modules/merchant/repository"
	merchantAuthService "go-backend/modules/merchant/service"

	// Payment Link
	paymentLinkController "go-backend/modules/payment_link/controller"
	paymentLinkRepo "go-backend/modules/payment_link/repository"
	paymentLinkService "go-backend/modules/payment_link/service"

	// Payout
	payoutController "go-backend/modules/payout/controller"
	payoutRepo "go-backend/modules/payout/repository"
	payoutService "go-backend/modules/payout/service"

	// Settlement
	settlementClient "go-backend/modules/settlement/client"
	settlementController "go-backend/modules/settlement/controller"
	settlementRepo "go-backend/modules/settlement/repository"
	settlementService "go-backend/modules/settlement/service"

	// Transaction
	txController "go-backend/modules/transaction/controller"
	txRepo "go-backend/modules/transaction/repository"
	txService "go-backend/modules/transaction/service"

	// Wallet
	walletController "go-backend/modules/wallet/controller"
	walletRepo "go-backend/modules/wallet/repository"
	walletService "go-backend/modules/wallet/service"

	// Webhook
	webhookRepo "go-backend/modules/webhook/repository"
	webhookService "go-backend/modules/webhook/service"
)

var (
	kipayDomainRegex   = regexp.MustCompile(`^https:\/\/(.*\.)?kipay\.xyz$`)
	vercelDomainRegex  = regexp.MustCompile(`^https:\/\/(.*\.)?vercel\.app$`)
	netlifyDomainRegex = regexp.MustCompile(`^https:\/\/(.*\.)?netlify\.app$`)
)

func getCORSHandler() *cors.Cors {
	return cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			// Local development
			switch origin {
			case
				"http://localhost:5173",
				"http://127.0.0.1:5173",
				"http://localhost:5500",
				"http://127.0.0.1:5500":
				return true
			}

			// Production domains
			if kipayDomainRegex.MatchString(origin) ||
				vercelDomainRegex.MatchString(origin) ||
				netlifyDomainRegex.MatchString(origin) {
				return true
			}

			return false
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
			"X-API-KEY",
			"X-Requested-With",
		},
		AllowCredentials: true,
	})
}

func main() {
	// ============================================================
	// 1. Database Connection
	// ============================================================
	db := database.ConnectDB()
	defer db.Close()

	// ============================================================
	// 2. Router & Base Path Setup (/api/v1)
	// ============================================================
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// ============================================================
	// 3. Middlewares
	// ============================================================
	authMiddleware := middleware.AuthMiddleware(db)
	apiKeyMiddleware := middleware.APIKeyMiddleware(db)

	// ============================================================
	// 4. Dependency Injection (Repositories, Services, Controllers)
	// ============================================================

	// Webhook
	webhookRepoInstance := webhookRepo.NewWebhookRepository(db)
	webhookServiceInstance := webhookService.NewWebhookService(webhookRepoInstance)

	// Merchant Authentication
	merchantAuthRepoInstance := merchantAuthRepo.NewMerchantRepository(db)
	merchantAuthServiceInstance := merchantAuthService.NewMerchantAuthService(merchantAuthRepoInstance)
	merchantAuthCtrl := merchantAuthController.NewMerchantAuthController(merchantAuthServiceInstance)

	// Invoice
	invoiceRepoInstance := invoiceRepo.NewInvoiceRepository(db)
	invoiceServiceInstance := invoiceService.NewInvoiceService(invoiceRepoInstance, merchantAuthRepoInstance)
	invoiceCtrl := invoiceController.NewInvoiceController(invoiceServiceInstance)

	// Payment Links
	paymentLinkRepoInstance := paymentLinkRepo.NewPaymentLinkRepository(db)
	paymentLinkServiceInstance := paymentLinkService.NewPaymentLinkService(paymentLinkRepoInstance)
	paymentLinkCtrl := paymentLinkController.NewPaymentLinkController(paymentLinkServiceInstance)

	// API Keys
	apiKeyRepoInstance := apiKeyRepo.NewApiKeyRepository(db)
	apiKeyServiceInstance := apiKeyService.NewApiKeyService(apiKeyRepoInstance)
	apiKeyCtrl := apiKeyController.NewApiKeyController(apiKeyServiceInstance)

	// KYC
	kycRepoInstance := kycRepo.NewKycRepository(db)
	kycServiceInstance := kycService.NewKycService(kycRepoInstance)
	kycCtrl := kycController.NewKycController(kycServiceInstance)

	// Admin
	adminRepoInstance := adminRepo.NewAdminRepository(db)
	adminServiceInstance := adminService.NewAdminService(adminRepoInstance)
	adminCtrl := adminController.NewAdminController(adminServiceInstance)

	// Dashboard
	dashboardRepoInstance := dashboardRepo.NewDashboardRepository(db)
	dashboardServiceInstance := dashboardService.NewDashboardService(dashboardRepoInstance)
	dashboardCtrl := dashboardController.NewDashboardController(dashboardServiceInstance)

	// Settlement / Rust gRPC
	grpcAddr := os.Getenv("RUST_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051"
	}
	rustClient, err := settlementClient.NewRustVerificationClient(grpcAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Rust verification engine via gRPC: %v", err)
	}
	log.Printf("Connected to Rust verification engine at %s", grpcAddr)

	settlementRepoInstance := settlementRepo.NewSettlementRepository(db)
	settlementServiceInstance := settlementService.NewSettlementService(settlementRepoInstance, webhookServiceInstance, rustClient)
	settlementCtrl := settlementController.NewSettlementController(settlementServiceInstance)

	// Payout
	payoutRepoInstance := payoutRepo.NewPayoutRepository(db)
	payoutServiceInstance := payoutService.NewPayoutService(payoutRepoInstance)
	payoutCtrl := payoutController.NewPayoutController(payoutServiceInstance)

	// Transactions
	txRepoInstance := txRepo.NewTransactionRepository(db)
	txServiceInstance := txService.NewTransactionService(txRepoInstance)
	txCtrl := txController.NewTransactionController(txServiceInstance)

	// Wallet
	walletRepoInstance := walletRepo.NewWalletRepository(db)
	walletServiceInstance := walletService.NewWalletService(walletRepoInstance)
	walletCtrl := walletController.NewWalletController(walletServiceInstance)

	// ============================================================
	// 5. ROUTE DEFINITIONS
	// ============================================================

	// ------------------------------------------------------------
	// A. PUBLIC ROUTES (No authentication required)
	// ------------------------------------------------------------
	apiRouter.HandleFunc("/auth/register", merchantAuthCtrl.Register).Methods(http.MethodPost)
	apiRouter.HandleFunc("/auth/login", merchantAuthCtrl.Login).Methods(http.MethodPost)
	apiRouter.HandleFunc("/auth/logout", merchantAuthCtrl.Logout).Methods(http.MethodPost)

	// Public Payment Links & Invoices
	apiRouter.HandleFunc("/payment-links/{linkId}", paymentLinkCtrl.GetPAyentLinkById).Methods(http.MethodGet)
	apiRouter.HandleFunc("/payment-links/invoices", invoiceCtrl.CreateLinkInvoice).Methods(http.MethodPost)
	apiRouter.HandleFunc("/invoices/{id}/status", invoiceCtrl.GetInvoiceStatus).Methods(http.MethodGet)

	// Internal Settlement Webhooks
	apiRouter.HandleFunc("/internal/settlement", settlementCtrl.HandleVerificationEvent).Methods(http.MethodPost)


	// ------------------------------------------------------------
	// B. EXTERNAL API ROUTES (Protected by API Key Middleware)
	// Base path: /api/v1/external/...
	// ------------------------------------------------------------
	apiApiKeyRouter := apiRouter.PathPrefix("/external").Subrouter()
	apiApiKeyRouter.Use(apiKeyMiddleware)

	// Accessible at: POST /api/v1/external/invoices
	apiApiKeyRouter.HandleFunc("/invoices", invoiceCtrl.CreateDirectInvoice).Methods(http.MethodPost)


	// ------------------------------------------------------------
	// C. MERCHANT DASHBOARD ROUTES (Protected by Session/Auth Middleware)
	// Base path: /api/v1/...
	// ------------------------------------------------------------
	protectedRouter := apiRouter.PathPrefix("").Subrouter()
	protectedRouter.Use(authMiddleware)

	// Merchant Profile
	protectedRouter.HandleFunc("/auth/me", merchantAuthCtrl.GetCurrentUser).Methods(http.MethodGet)

	// Payment Links
	protectedRouter.HandleFunc("/payment-links", paymentLinkCtrl.Create).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/payment-links", paymentLinkCtrl.GetAllPaymentLinkByMerchant).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/transaction/invoices/{linkId}", paymentLinkCtrl.Update).Methods(http.MethodPut)

	// Invoices
	protectedRouter.HandleFunc("/link/invoices", invoiceCtrl.CreateLinkInvoice).Methods(http.MethodPost)

	// API Keys Management
	protectedRouter.HandleFunc("/api-keys", apiKeyCtrl.Create).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/api-keys", apiKeyCtrl.GetAll).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/api-keys/{id}", apiKeyCtrl.Delete).Methods(http.MethodDelete)

	// KYC
	protectedRouter.HandleFunc("/kyc/documents", kycCtrl.Submit).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/kyc/status", kycCtrl.GetStatus).Methods(http.MethodGet)

	// Dashboard Metrics
	protectedRouter.HandleFunc("/dashboard/metrics", dashboardCtrl.GetMetrics).Methods(http.MethodGet)

	// Payouts
	protectedRouter.HandleFunc("/payouts", payoutCtrl.Create).Methods(http.MethodPost)

	// Wallets Management
	protectedRouter.HandleFunc("/wallets", walletCtrl.Create).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/wallets", walletCtrl.List).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/wallets/{id}", walletCtrl.GetByID).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/wallets/{id}", walletCtrl.Update).Methods(http.MethodPut)
	protectedRouter.HandleFunc("/wallets/{id}", walletCtrl.Delete).Methods(http.MethodDelete)

	// Transactions
	protectedRouter.HandleFunc("/transactions/invoices", txCtrl.ListInvoices).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/transactions/link-invoices", txCtrl.ListAllLinkInvoices).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/transaction/invoices/{payment-link-id}", txCtrl.ListAllLinkInvoicesByPaymentLinkId).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/transactions/settled", txCtrl.ListTransactions).Methods(http.MethodGet)


	// ------------------------------------------------------------
	// D. ADMIN ROUTES
	// ------------------------------------------------------------
	apiRouter.HandleFunc("/admin/merchants", adminCtrl.ListMerchants).Methods(http.MethodGet)
	apiRouter.HandleFunc("/admin/merchants/{merchantId}/status", adminCtrl.UpdateMerchantStatus).Methods(http.MethodPut)


	// ============================================================
	// 6. Server Initialization & CORS
	// ============================================================
	corsHandler := getCORSHandler()
	handler := corsHandler.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("kipay.xyz backend running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}