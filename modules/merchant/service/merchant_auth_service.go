package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"go-backend/modules/merchant/dto"
	"go-backend/modules/merchant/repository"
)

type MerchantAuthService struct {
	repo *repository.MerchantRepository
}

func NewMerchantAuthService(repo *repository.MerchantRepository) *MerchantAuthService {
	return &MerchantAuthService{repo: repo}
}

func (s *MerchantAuthService) Register(req dto.RegisterMerchantRequest) (*dto.MerchantAuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	merchantID := fmt.Sprintf("mch_%d", time.Now().UnixNano())
	merchant, err := s.repo.Create(merchantID, req.BusinessName, req.Email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("email already registered or DB error: %v", err)
	}

	token, err := generateJWT(merchant.ID, merchant.Email)
	if err != nil {
		return nil, err
	}

	return &dto.MerchantAuthResponse{
		Token:        token,
		MerchantID:   merchant.ID,
		BusinessName: merchant.BusinessName,
		Email:        merchant.Email,
		Status:       merchant.Status,
	}, nil
}

func (s *MerchantAuthService) Login(req dto.LoginMerchantRequest) (*dto.MerchantAuthResponse, error) {
	merchant, passwordHash, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := generateJWT(merchant.ID, merchant.Email)
	if err != nil {
		return nil, err
	}

	return &dto.MerchantAuthResponse{
		Token:        token,
		MerchantID:   merchant.ID,
		BusinessName: merchant.BusinessName,
		Email:        merchant.Email,
		Status:       merchant.Status,
	}, nil
}

func generateJWT(merchantID string, email string) (string, error) {
	claims := jwt.MapClaims{
		"merchant_id": merchantID,
		"email":       email,
		"exp":         time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "go-backend_super_secret_key"
	}
	return token.SignedString([]byte(secret))
}
func (s *MerchantAuthService) GetMeService(merchantID string) (*dto.GetMeResponse, error) {
    // Call GetByID using the passed merchant_id directly
    merchant, _, err := s.repo.GetByID(merchantID)
    if err != nil {
        return nil, errors.New("merchant profile not found")
    }

    // Format timestamps and nullable fields safely
    var kycSubmitted, verifiedAt *string
    if merchant.KycSubmittedAt != nil {
        t := merchant.KycSubmittedAt.Format(time.RFC3339)
        kycSubmitted = &t
    }
    if merchant.VerifiedAt != nil {
        t := merchant.VerifiedAt.Format(time.RFC3339)
        verifiedAt = &t
    }

    return &dto.GetMeResponse{
        ID:               merchant.ID,
        BusinessName:     merchant.BusinessName,
        Email:            merchant.Email,
        Status:           merchant.Status,
        SolanaWallet:     merchant.SolanaWallet,
        PolygonWallet:    merchant.PolygonWallet,
        EthereumWallet:   merchant.EthereumWallet,
        KycSubmittedAt:   kycSubmitted,
        VerifiedAt:       verifiedAt,
        CreatedAt:        merchant.CreatedAt.Format(time.RFC3339),
    }, nil
}
