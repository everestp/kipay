package service

import (
    "errors"
    "fmt"
    "time"

    "go-backend/modules/invoice/dto"
    "go-backend/modules/invoice/repository"
	 merchantRepo "go-backend/modules/merchant/repository"
	 mo "go-backend/models"
)


type InvoiceService struct {
    invoiceRepo *repository.InvoiceRepository
    merchantRepo *merchantRepo.MerchantRepository // Needed to look up merchant wallet addresses & payment link owners
}

func NewInvoiceService(invoiceRepo *repository.InvoiceRepository, merchantRepo *merchantRepo.MerchantRepository) *InvoiceService {
    return &InvoiceService{
        invoiceRepo:  invoiceRepo,
        merchantRepo: merchantRepo,
    }
}

// ==========================================
// 1. CREATE PAYMENT LINK INVOICE (Public Flow)
// ==========================================
func (s *InvoiceService) CreateLinkInvoice(req dto.CreateLinkInvoiceRequest) (*dto.InvoiceResponse, error) {
    // 1. Fetch payment link details from DB to verify it exists and get its merchant_id
    link, err := s.invoiceRepo.GetPaymentLinkByID(req.PaymentLinkId)
    if err != nil || !link.IsActive {
        return nil, errors.New("invalid or inactive payment link")
    }

    // 2. Fetch merchant's specific multi-chain wallet address
    merchantWallet, err := s.merchantRepo.GetWalletByNetwork(link.MerchantID, req.Network)
    if err != nil || merchantWallet == "" {
        return nil, fmt.Errorf("merchant has not configured a payout wallet for network: %s", req.Network)
    }

    // 3. Calculate crypto amount based on live/mock rate
    cryptoRate := getMockCryptoRate(req.Currency)
    if cryptoRate <= 0 {
        return nil, errors.New("unsupported currency")
    }

    // If payment link allows custom amounts, use req.AmountUSD, otherwise use link's fixed amount
    amountUSD := link.AmountUSD
    if link.PricingType == "CUSTOM" {
        if req.AmountUSD < link.MinAmountUSD {
            return nil, fmt.Errorf("amount is less than minimum allowed ($%.2f)", link.MinAmountUSD)
        }
        amountUSD = req.AmountUSD
    }

    amountCrypto := amountUSD / cryptoRate
    invoiceID := fmt.Sprintf("inv_pl_%d", time.Now().UnixNano())
    expiresAt := time.Now().Add(time.Minute * 15) // 15 mins expiry window

    // 4. Generate QR code payload data based on network
    qrCodeData := generateQRCodeURI(req.Network, merchantWallet, amountCrypto, invoiceID)

    // 5. Save into payment_link_invoices table
    inv, err := s.invoiceRepo.CreatePaymentLinkInvoice(
        invoiceID, link.MerchantID, link.ID, amountUSD, req.Currency,
        req.Network, amountCrypto, merchantWallet, qrCodeData, expiresAt,
    )
    if err != nil {
        return nil, err
    }

    return mapToInvoiceResponse(inv), nil
}

// ==========================================
// 2. CREATE DIRECT INVOICE (E-Commerce API Flow)
// ==========================================
func (s *InvoiceService) CreateDirectInvoice(merchantID string, req dto.CreateDirectInvoiceRequest) (*dto.InvoiceResponse, error) {
    // 1. Fetch merchant's specific multi-chain wallet address
    merchantWallet, err := s.merchantRepo.GetWalletByNetwork(merchantID, req.Network)
    if err != nil || merchantWallet == "" {
        return nil, fmt.Errorf("merchant has not configured a payout wallet for network: %s", req.Network)
    }

    // 2. Calculate crypto amount
    cryptoRate := getMockCryptoRate(req.Currency)
    if cryptoRate <= 0 {
        return nil, errors.New("unsupported currency")
    }

    amountCrypto := req.AmountUSD / cryptoRate
    invoiceID := fmt.Sprintf("inv_dir_%d", time.Now().UnixNano())
    expiresAt := time.Now().Add(time.Minute * 15)

    // 3. Generate QR code payload
    qrCodeData := generateQRCodeURI(req.Network, merchantWallet, amountCrypto, invoiceID)

    // 4. Save into direct_invoices table (with order_id)
    inv, err := s.invoiceRepo.CreateDirectInvoice(
        invoiceID, merchantID, req.OrderID, req.AmountUSD, req.Currency,
        req.Network, amountCrypto, merchantWallet, qrCodeData, expiresAt,
    )
    if err != nil {
        return nil, err
    }

    return mapToInvoiceResponse(inv), nil
}

// ==========================================
// 3. GET INVOICE STATUS (Public Polling)
// ==========================================
func (s *InvoiceService) GetInvoiceStatus(invoiceID string) (*dto.InvoiceResponse, error) {
    // Looks up across both invoice tables or handles it based on ID prefix
    inv, err := s.invoiceRepo.FindInvoiceAnyType(invoiceID)
    if err != nil {
        return nil, errors.New("invoice not found")
    }

    return mapToInvoiceResponse(inv), nil
}

// ==========================================
// HELPER UTILITIES
// ==========================================
func getMockCryptoRate(currency string) float64 {
    switch currency {
    case "USDT", "USDC":
        return 1.00
    case "SOL":
        return 180.50
    case "ETH":
        return 2850.00
    case "POL":
        return 0.55
    default:
        return 0.00
    }
}

func generateQRCodeURI(network, address string, amountCrypto float64, invoiceID string) string {
    switch network {
    case "solana":
        return fmt.Sprintf("solana:%s?amount=%.8f&message=Invoice%%20%s", address, amountCrypto, invoiceID)
    case "polygon", "ethereum":
        return fmt.Sprintf("ethereum:%s@%s?value=%.8f", address, network, amountCrypto)
    default:
        return fmt.Sprintf("%s:%s?amount=%.8f", network, address, amountCrypto)
    }
}

func mapToInvoiceResponse(inv *mo.Invoice) *dto.InvoiceResponse {
    confirmedAtStr := ""
    if inv.ConfirmedAt != nil {
        confirmedAtStr = inv.ConfirmedAt.Format(time.RFC3339)
    }

    return &dto.InvoiceResponse{
        InvoiceID:      inv.ID,
        OrderID:         inv.OrderID,
        MerchantID:     inv.MerchantID,
        AmountUSD:      inv.AmountUSD,
        Currency:       inv.Currency,
        Network:        inv.Network,
        AmountCrypto:   inv.AmountCrypto,
        Status:         inv.Status,
        DepositAddress: inv.DepositAddress,
        QRCodeData:     inv.QRCodeData,
        ExpiresAt:      inv.ExpiresAt.Format(time.RFC3339),
        ConfirmedAt:    confirmedAtStr,
    }
}
// ==========================================
// GET DIRECT INVOICE BY ID
// ==========================================
func (s *InvoiceService) GetDirectInvoiceByID(invoiceID string) (*dto.InvoiceResponse, error) {
	inv, err := s.invoiceRepo.GetDirectInvoiceByID(invoiceID)
	if err != nil {
		return nil, errors.New("direct invoice not found")
	}

	return mapToInvoiceResponse(inv), nil
}
