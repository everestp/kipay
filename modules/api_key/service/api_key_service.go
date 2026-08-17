package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"go-backend/modules/api_key/dto"
	"go-backend/modules/api_key/repository"
)

type ApiKeyService struct {
	repo *repository.ApiKeyRepository
}

func NewApiKeyService(repo *repository.ApiKeyRepository) *ApiKeyService {
	return &ApiKeyService{repo: repo}
}

func (s *ApiKeyService) GenerateApiKey(merchantID string, req dto.CreateApiKeyRequest) (*dto.ApiKeyResponse, error) {
	// Generate random bytes for prefix and secret
	prefixBytes := make([]byte, 8)
	rand.Read(prefixBytes)
	keyPrefix := fmt.Sprintf("pine_live_%s", hex.EncodeToString(prefixBytes))

	secretBytes := make([]byte, 24)
	rand.Read(secretBytes)
	rawSecret := hex.EncodeToString(secretBytes)

	// Full key format displayed once: pine_live_..._<rawSecret>
	fullKeyDisplay := fmt.Sprintf("%s_%s", keyPrefix, rawSecret)

	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(rawSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("key_%d", time.Now().UnixNano())
	res, err := s.repo.Create(id, merchantID, req.Name, keyPrefix, string(hashedSecret))
	if err != nil {
		return nil, err
	}

	// Attach the unhashed secret ONLY on initial generation response
	res.SecretKey = fullKeyDisplay
	return res, nil
}

func (s *ApiKeyService) DeleteApiKey(merchantID string, keyID string) error {
	if err := s.repo.Delete(keyID, merchantID); err != nil {
		return fmt.Errorf("failed to delete api key: %w", err)
	}

	return nil
}
func (s *ApiKeyService) GetApiKeys(merchantID string) ([]dto.ApiKeyResponse, error) {
	return s.repo.GetAllByMerchant(merchantID)
}
