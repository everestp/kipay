package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"time"

	"go-backend/modules/webhook/dto"
	"go-backend/modules/webhook/repository"
)

type WebhookService struct {
	repo *repository.WebhookRepository
}

func NewWebhookService(repo *repository.WebhookRepository) *WebhookService {
	return &WebhookService{repo: repo}
}

func (s *WebhookService) RegisterEndpoint(merchantID string, req dto.RegisterWebhookEndpointRequest) error {
	id := fmt.Sprintf("wh_ep_%d", time.Now().UnixNano())

	secretBytes := make([]byte, 16)
	rand.Read(secretBytes)
	secret := hex.EncodeToString(secretBytes)

	return s.repo.SaveEndpoint(id, merchantID, req.EndpointURL, req.SubscribedEvents, secret)
}

func (s *WebhookService) DispatchEvent(merchantID string, eventType string, data interface{}) error {
	endpoints, err := s.repo.GetActiveEndpointsForMerchant(merchantID, eventType)
	if err != nil || len(endpoints) == 0 {
		return nil // No active webhooks for this event
	}

	payload := dto.WebhookEventPayload{
		EventID:   fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		EventType: eventType,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, ep := range endpoints {
		// Compute HMAC-SHA256 signature
		mac := hmac.New(sha256.New, []byte(ep.Secret))
		mac.Write(bodyBytes)
		signature := hex.EncodeToString(mac.Sum(nil))

		req, err := http.NewRequest("POST", ep.URL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Pinecone-Signature", fmt.Sprintf("sha256=%s", signature))

		resp, err := client.Do(req)
		status := "SUCCESS"
		respCode := 0
		if err != nil || resp.StatusCode >= 400 {
			status = "FAILED"
			if resp != nil {
				respCode = resp.StatusCode
			}
		} else {
			respCode = resp.StatusCode
		}
		if resp != nil {
			resp.Body.Close()
		}

		logID := fmt.Sprintf("wh_log_%d", time.Now().UnixNano())
		_ = s.repo.LogDelivery(logID, merchantID, eventType, ep.URL, respCode, status)
	}

	return nil
}
