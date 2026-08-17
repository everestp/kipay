package repository

import (
	"database/sql"

	"strings"


	
)

type WebhookRepository struct {
	db *sql.DB
}

func NewWebhookRepository(db *sql.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) SaveEndpoint(id string, merchantID string, url string, events []string, secret string) error {
	query := `
		INSERT INTO webhook_endpoints (id, merchant_id, endpoint_url, subscribed_events, secret, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, TRUE, NOW())
		ON CONFLICT (id) DO UPDATE
		SET endpoint_url = EXCLUDED.endpoint_url, subscribed_events = EXCLUDED.subscribed_events
	`
	eventArr := "{" + strings.Join(events, ",") + "}"
	_, err := r.db.Exec(query, id, merchantID, url, eventArr, secret)
	return err
}

func (r *WebhookRepository) GetActiveEndpointsForMerchant(merchantID string, eventType string) ([]struct{ URL, Secret string }, error) {
	query := `
		SELECT endpoint_url, secret
		FROM webhook_endpoints
		WHERE merchant_id = $1 AND is_active = TRUE AND $2 = ANY(subscribed_events)
	`
	rows, err := r.db.Query(query, merchantID, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []struct{ URL, Secret string }
	for rows.Next() {
		var ep struct{ URL, Secret string }
		if err := rows.Scan(&ep.URL, &ep.Secret); err == nil {
			endpoints = append(endpoints, ep)
		}
	}
	return endpoints, nil
}

func (r *WebhookRepository) LogDelivery(logID string, merchantID string, eventType string, endpointURL string, responseCode int, status string) error {
	query := `
		INSERT INTO webhook_logs (id, merchant_id, event_type, endpoint_url, response_code, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err := r.db.Exec(query, logID, merchantID, eventType, endpointURL, responseCode, status)
	return err
}
