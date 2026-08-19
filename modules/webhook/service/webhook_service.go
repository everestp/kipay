package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-backend/modules/webhook/dto"
	"go-backend/modules/webhook/repository"
)

const (
	SignatureHeader = "X-Kipay-Signature"
	SignaturePrefix = "sha256="
)

type WebhookService struct {
	repo *repository.WebhookRepository
	http *http.Client
}

func NewWebhookService(
	repo *repository.WebhookRepository,
) *WebhookService {
	return &WebhookService{
		repo: repo,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ============================================================
// REGISTER WEBHOOK
// ============================================================

func (s *WebhookService) RegisterEndpoint(
	merchantID string,
	req dto.RegisterWebhookEndpointRequest,
) (string, string, error) {

	if merchantID == "" {
		return "", "", fmt.Errorf(
			"merchant ID is required",
		)
	}

	if req.EndpointURL == "" {
		return "", "", fmt.Errorf(
			"endpoint URL is required",
		)
	}

	if len(req.SubscribedEvents) == 0 {
		return "", "", fmt.Errorf(
			"at least one event must be subscribed",
		)
	}

	/*
	 * Generate webhook endpoint ID.
	 */
	id, err := generateID("wh_ep")
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to generate webhook ID: %w",
			err,
		)
	}

	/*
	 * Generate cryptographically secure secret.
	 */
	secret, err := generateSecret()
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to generate webhook secret: %w",
			err,
		)
	}

	/*
	 * Save endpoint.
	 */
	err = s.repo.SaveEndpoint(
		id,
		merchantID,
		req.EndpointURL,
		req.SubscribedEvents,
		secret,
	)

	if err != nil {
		return "", "", err
	}

	/*
	 * IMPORTANT:
	 *
	 * Secret is returned only during creation.
	 *
	 * Your controller should return this to the merchant
	 * exactly once.
	 */
	return id, secret, nil
}

// ============================================================
// LIST WEBHOOK ENDPOINTS
// ============================================================

func (s *WebhookService) ListEndpoints(
	merchantID string,
) ([]repository.WebhookEndpointInfo, error) {

	if merchantID == "" {
		return nil, fmt.Errorf(
			"merchant ID is required",
		)
	}

	endpoints, err :=
		s.repo.ListEndpoints(merchantID)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to list webhook endpoints: %w",
			err,
		)
	}

	return endpoints, nil
}

// ============================================================
// ROTATE WEBHOOK SECRET
// ============================================================

func (s *WebhookService) RotateSecret(
	merchantID string,
	endpointID string,
) (string, error) {

	if merchantID == "" {
		return "", fmt.Errorf(
			"merchant ID is required",
		)
	}

	if endpointID == "" {
		return "", fmt.Errorf(
			"webhook endpoint ID is required",
		)
	}

	/*
	 * Verify that this endpoint belongs
	 * to the authenticated merchant.
	 */
	_, err := s.repo.GetEndpoint(
		merchantID,
		endpointID,
	)

	if err != nil {
		return "", fmt.Errorf(
			"failed to find webhook endpoint: %w",
			err,
		)
	}

	/*
	 * Generate completely new secret.
	 */
	newSecret, err := generateSecret()

	if err != nil {
		return "", fmt.Errorf(
			"failed to generate new webhook secret: %w",
			err,
		)
	}

	/*
	 * Update secret in database.
	 */
	err = s.repo.UpdateSecret(
		merchantID,
		endpointID,
		newSecret,
	)

	if err != nil {
		return "", fmt.Errorf(
			"failed to rotate webhook secret: %w",
			err,
		)
	}

	/*
	 * Return new secret ONCE.
	 */
	return newSecret, nil
}

// ============================================================
// DEACTIVATE WEBHOOK
// ============================================================

func (s *WebhookService) DeactivateEndpoint(
	merchantID string,
	endpointID string,
) error {

	if merchantID == "" {
		return fmt.Errorf(
			"merchant ID is required",
		)
	}

	if endpointID == "" {
		return fmt.Errorf(
			"webhook endpoint ID is required",
		)
	}

	err := s.repo.DeactivateEndpoint(
		merchantID,
		endpointID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to deactivate webhook: %w",
			err,
		)
	}

	return nil
}

// ============================================================
// DISPATCH EVENT
// ============================================================

func (s *WebhookService) DispatchEvent(
	merchantID string,
	eventType string,
	data interface{},
) error {

	if merchantID == "" {
		return fmt.Errorf(
			"merchant ID is required",
		)
	}

	if eventType == "" {
		return fmt.Errorf(
			"event type is required",
		)
	}

	/*
	 * Find all active endpoints subscribed
	 * to this event.
	 */
	endpoints, err :=
		s.repo.GetActiveEndpointsForMerchant(
			merchantID,
			eventType,
		)

	if err != nil {
		return fmt.Errorf(
			"failed to get webhook endpoints: %w",
			err,
		)
	}

	/*
	 * No webhook configured.
	 */
	if len(endpoints) == 0 {
		return nil
	}

	/*
	 * Generate event ID.
	 */
	eventID, err := generateID("evt")

	if err != nil {
		return fmt.Errorf(
			"failed to generate event ID: %w",
			err,
		)
	}

	/*
	 * Build payload.
	 */
	payload := dto.WebhookEventPayload{
		EventID:   eventID,
		EventType: eventType,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}

	/*
	 * Marshal once.
	 *
	 * The exact bytes are signed and sent.
	 */
	bodyBytes, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf(
			"failed to marshal webhook payload: %w",
			err,
		)
	}

	/*
	 * Deliver to every subscribed endpoint.
	 *
	 * One endpoint failing must NOT stop
	 * the other endpoints.
	 */
	for _, endpoint := range endpoints {

		err := s.deliver(
			merchantID,
			eventType,
			eventID,
			endpoint.URL,
			endpoint.Secret,
			bodyBytes,
		)

		if err != nil {
			continue
		}
	}

	return nil
}

// ============================================================
// DELIVER WEBHOOK
// ============================================================

func (s *WebhookService) deliver(
	merchantID string,
	eventType string,
	eventID string,
	endpointURL string,
	secret string,
	bodyBytes []byte,
) error {

	/*
	 * Generate HMAC-SHA256 signature.
	 */
	signature := generateSignature(
		secret,
		bodyBytes,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		endpointURL,
		bytes.NewReader(bodyBytes),
	)

	if err != nil {

		s.logDelivery(
			merchantID,
			eventType,
			eventID,
			endpointURL,
			0,
			"FAILED",
		)

		return err
	}

	/*
	 * Standard headers.
	 */
	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"Accept",
		"application/json",
	)

	req.Header.Set(
		"User-Agent",
		"Kipay-Webhooks/1.0",
	)

	/*
	 * Signature.
	 */
	req.Header.Set(
		SignatureHeader,
		SignaturePrefix+signature,
	)

	/*
	 * Event metadata.
	 */
	req.Header.Set(
		"X-Kipay-Event",
		eventType,
	)

	req.Header.Set(
		"X-Kipay-Event-ID",
		eventID,
	)

	/*
	 * Send request.
	 */
	resp, err := s.http.Do(req)

	if err != nil {

		s.logDelivery(
			merchantID,
			eventType,
			eventID,
			endpointURL,
			0,
			"FAILED",
		)

		return err
	}

	defer resp.Body.Close()

	/*
	 * Consume response body so the HTTP connection
	 * can be reused.
	 */
	_, _ = io.Copy(
		io.Discard,
		resp.Body,
	)

	/*
	 * Any 2xx response means successful delivery.
	 */
	if resp.StatusCode >= 200 &&
		resp.StatusCode < 300 {

		s.logDelivery(
			merchantID,
			eventType,
			eventID,
			endpointURL,
			resp.StatusCode,
			"SUCCESS",
		)

		return nil
	}

	/*
	 * Non-2xx response.
	 */
	s.logDelivery(
		merchantID,
		eventType,
		eventID,
		endpointURL,
		resp.StatusCode,
		"FAILED",
	)

	return fmt.Errorf(
		"webhook returned HTTP %d",
		resp.StatusCode,
	)
}

// ============================================================
// GENERATE HMAC SIGNATURE
// ============================================================

func generateSignature(
	secret string,
	body []byte,
) string {

	mac := hmac.New(
		sha256.New,
		[]byte(secret),
	)

	_, _ = mac.Write(body)

	return hex.EncodeToString(
		mac.Sum(nil),
	)
}

// ============================================================
// GENERATE WEBHOOK SECRET
// ============================================================

func generateSecret() (string, error) {

	/*
	 * 32 random bytes = 256 bits.
	 */
	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return "whsec_" +
		hex.EncodeToString(buf), nil
}

// ============================================================
// GENERATE ID
// ============================================================

func generateID(
	prefix string,
) (string, error) {

	buf := make([]byte, 12)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s_%s",
		prefix,
		hex.EncodeToString(buf),
	), nil
}

// ============================================================
// LOG DELIVERY
// ============================================================

func (s *WebhookService) logDelivery(
	merchantID string,
	eventType string,
	eventID string,
	endpointURL string,
	responseCode int,
	status string,
) {

	/*
	 * Webhook logging must NEVER break
	 * payment settlement.
	 */
	logID, err := generateID("wh_log")

	if err != nil {
		return
	}

	_ = s.repo.LogDelivery(
		logID,
		merchantID,
		eventType,
		endpointURL,
		responseCode,
		status,
	)
}
type TestWebhookPayload struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	LiveMode  bool                   `json:"live_mode"`
	CreatedAt string                 `json:"created_at"`
	Data      map[string]interface{} `json:"data"`
}

func (s *WebhookService) SendTestWebhook(
	ctx context.Context,
	webhookID string,
	merchantID string,
) error {

	webhook, err :=
		s.repo.GetWebhookEndpoint(
			ctx,
			webhookID,
			merchantID,
		)

	if err != nil {
		return err
	}

	if !webhook.IsActive {
		return fmt.Errorf(
			"webhook endpoint is inactive",
		)
	}

	payload := TestWebhookPayload{
		ID:        fmt.Sprintf(
			"evt_test_%d",
			time.Now().UnixNano(),
		),
		Type:      "invoice.confirmed",
		LiveMode:  false,
		CreatedAt: time.Now().UTC().Format(
			time.RFC3339,
		),
		Data: map[string]interface{}{
			"test": true,
			"invoice": map[string]interface{}{
				"id":     "inv_test_123456",
				"status": "confirmed",
			},
			"payment": map[string]interface{}{
				"tx_hash": "test_transaction_hash",
				"network": "testnet",
				"currency": "USDC",
				"amount":   10.00,
			},
		},
	}

	body, err :=
		json.Marshal(payload)

	if err != nil {
		return fmt.Errorf(
			"failed to encode test webhook payload: %w",
			err,
		)
	}

	req, err :=
		http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			webhook.URL,
			bytes.NewReader(body),
		)

	if err != nil {
		return fmt.Errorf(
			"failed to create webhook request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"User-Agent",
		"Kipay-Webhook-Test/1.0",
	)

	req.Header.Set(
		"X-Kipay-Test",
		"true",
	)

	req.Header.Set(
		"X-Kipay-Event",
		"invoice.confirmed",
	)

	resp, err :=
		s.http.Do(req)

	if err != nil {
		return fmt.Errorf(
			"failed to send test webhook: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		return fmt.Errorf(
			"webhook endpoint returned status %d",
			resp.StatusCode,
		)
	}

	return nil
}
