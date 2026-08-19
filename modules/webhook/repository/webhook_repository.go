package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type WebhookEndpoint struct {
	URL    string
	Secret string
}

type WebhookRepository struct {
	db *sql.DB
}

func NewWebhookRepository(db *sql.DB) *WebhookRepository {
	return &WebhookRepository{
		db: db,
	}
}

// ============================================================
// SAVE / UPDATE WEBHOOK ENDPOINT
// ============================================================

func (r *WebhookRepository) SaveEndpoint(
	id string,
	merchantID string,
	url string,
	events []string,
	secret string,
) error {

	if r.db == nil {
		return errors.New("webhook repository: database is nil")
	}

	if id == "" {
		return errors.New("webhook repository: endpoint id is required")
	}

	if merchantID == "" {
		return errors.New("webhook repository: merchant id is required")
	}

	if url == "" {
		return errors.New("webhook repository: endpoint URL is required")
	}

	if len(events) == 0 {
		return errors.New("webhook repository: at least one event is required")
	}

	if secret == "" {
		return errors.New("webhook repository: webhook secret is required")
	}

	const query = `
		INSERT INTO webhook_endpoints (
			id,
			merchant_id,
			endpoint_url,
			subscribed_events,
			secret,
			is_active,
			created_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			TRUE,
			NOW()
		)
		ON CONFLICT (id)
		DO UPDATE SET
			endpoint_url = EXCLUDED.endpoint_url,
			subscribed_events = EXCLUDED.subscribed_events,
			is_active = TRUE
	`

	_, err := r.db.Exec(
		query,
		id,
		merchantID,
		url,
		pq.Array(events),
		secret,
	)

	if err != nil {
		return fmt.Errorf(
			"webhook repository: save endpoint: %w",
			err,
		)
	}

	return nil
}

// ============================================================
// GET ACTIVE ENDPOINTS FOR MERCHANT
// ============================================================

func (r *WebhookRepository) GetActiveEndpointsForMerchant(
	merchantID string,
	eventType string,
) ([]WebhookEndpoint, error) {

	if r.db == nil {
		return nil, errors.New(
			"webhook repository: database is nil",
		)
	}

	if merchantID == "" {
		return nil, errors.New(
			"webhook repository: merchant id is required",
		)
	}

	if eventType == "" {
		return nil, errors.New(
			"webhook repository: event type is required",
		)
	}

	const query = `
		SELECT
			endpoint_url,
			secret
		FROM webhook_endpoints
		WHERE
			merchant_id = $1
			AND is_active = TRUE
			AND $2 = ANY(subscribed_events)
	`

	rows, err := r.db.Query(
		query,
		merchantID,
		eventType,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"webhook repository: get active endpoints: %w",
			err,
		)
	}

	defer rows.Close()

	endpoints := make(
		[]WebhookEndpoint,
		0,
	)

	for rows.Next() {

		var endpoint WebhookEndpoint

		if err := rows.Scan(
			&endpoint.URL,
			&endpoint.Secret,
		); err != nil {

			return nil, fmt.Errorf(
				"webhook repository: scan endpoint: %w",
				err,
			)
		}

		endpoints = append(
			endpoints,
			endpoint,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"webhook repository: iterate endpoints: %w",
			err,
		)
	}

	return endpoints, nil
}

// ============================================================
// LOG WEBHOOK DELIVERY
// ============================================================

func (r *WebhookRepository) LogDelivery(
	logID string,
	merchantID string,
	eventType string,
	endpointURL string,
	responseCode int,
	status string,
) error {

	if r.db == nil {
		return errors.New(
			"webhook repository: database is nil",
		)
	}

	if logID == "" {
		return errors.New(
			"webhook repository: log id is required",
		)
	}

	if merchantID == "" {
		return errors.New(
			"webhook repository: merchant id is required",
		)
	}

	if eventType == "" {
		return errors.New(
			"webhook repository: event type is required",
		)
	}

	if endpointURL == "" {
		return errors.New(
			"webhook repository: endpoint URL is required",
		)
	}

	const query = `
		INSERT INTO webhook_logs (
			id,
			merchant_id,
			event_type,
			endpoint_url,
			response_code,
			status,
			created_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			NOW()
		)
	`

	_, err := r.db.Exec(
		query,
		logID,
		merchantID,
		eventType,
		endpointURL,
		responseCode,
		status,
	)

	if err != nil {
		return fmt.Errorf(
			"webhook repository: log delivery: %w",
			err,
		)
	}

	return nil
}

// ============================================================
// DEACTIVATE ENDPOINT
// ============================================================

func (r *WebhookRepository) DeactivateEndpoint(
	merchantID string,
	endpointID string,
) error {

	if r.db == nil {
		return errors.New(
			"webhook repository: database is nil",
		)
	}

	if merchantID == "" {
		return errors.New(
			"webhook repository: merchant id is required",
		)
	}

	if endpointID == "" {
		return errors.New(
			"webhook repository: endpoint id is required",
		)
	}

	const query = `
		UPDATE webhook_endpoints
		SET is_active = FALSE
		WHERE id = $1
		  AND merchant_id = $2
	`

	result, err := r.db.Exec(
		query,
		endpointID,
		merchantID,
	)

	if err != nil {
		return fmt.Errorf(
			"webhook repository: deactivate endpoint: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"webhook repository: rows affected: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ============================================================
// GET ENDPOINT BY ID
// ============================================================

func (r *WebhookRepository) GetEndpoint(
	merchantID string,
	endpointID string,
) (*WebhookEndpoint, error) {

	if r.db == nil {
		return nil, errors.New(
			"webhook repository: database is nil",
		)
	}

	const query = `
		SELECT
			endpoint_url,
			secret
		FROM webhook_endpoints
		WHERE id = $1
		  AND merchant_id = $2
	`

	var endpoint WebhookEndpoint

	err := r.db.QueryRow(
		query,
		endpointID,
		merchantID,
	).Scan(
		&endpoint.URL,
		&endpoint.Secret,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, fmt.Errorf(
			"webhook repository: get endpoint: %w",
			err,
		)
	}

	return &endpoint, nil
}

// ============================================================
// UPDATE WEBHOOK SECRET
// ============================================================

func (r *WebhookRepository) UpdateSecret(
	merchantID string,
	endpointID string,
	newSecret string,
) error {

	if r.db == nil {
		return errors.New(
			"webhook repository: database is nil",
		)
	}

	if merchantID == "" {
		return errors.New(
			"webhook repository: merchant id is required",
		)
	}

	if endpointID == "" {
		return errors.New(
			"webhook repository: endpoint id is required",
		)
	}

	if newSecret == "" {
		return errors.New(
			"webhook repository: new webhook secret is required",
		)
	}

	const query = `
		UPDATE webhook_endpoints
		SET
			secret = $1,
			is_active = TRUE
		WHERE
			id = $2
			AND merchant_id = $3
	`

	result, err := r.db.Exec(
		query,
		newSecret,
		endpointID,
		merchantID,
	)

	if err != nil {
		return fmt.Errorf(
			"webhook repository: update secret: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"webhook repository: rows affected: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
// ============================================================
// LIST MERCHANT WEBHOOK ENDPOINTS
// ============================================================

type WebhookEndpointInfo struct {
	ID               string
	URL              string
	SubscribedEvents []string
	IsActive         bool
	CreatedAt        string
}

func (r *WebhookRepository) ListEndpoints(
	merchantID string,
) ([]WebhookEndpointInfo, error) {

	if r.db == nil {
		return nil, errors.New(
			"webhook repository: database is nil",
		)
	}

	if merchantID == "" {
		return nil, errors.New(
			"webhook repository: merchant id is required",
		)
	}

	const query = `
		SELECT
			id,
			endpoint_url,
			subscribed_events,
			is_active,
			created_at
		FROM webhook_endpoints
		WHERE merchant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, merchantID)
	if err != nil {
		return nil, fmt.Errorf(
			"webhook repository: list endpoints: %w",
			err,
		)
	}

	defer rows.Close()

	endpoints := make(
		[]WebhookEndpointInfo,
		0,
	)

	for rows.Next() {

		var endpoint WebhookEndpointInfo

		var events pq.StringArray

		err := rows.Scan(
			&endpoint.ID,
			&endpoint.URL,
			&events,
			&endpoint.IsActive,
			&endpoint.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"webhook repository: scan endpoint: %w",
				err,
			)
		}

		endpoint.SubscribedEvents =
			[]string(events)

		endpoints = append(
			endpoints,
			endpoint,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"webhook repository: iterate endpoints: %w",
			err,
		)
	}

	return endpoints, nil
}

type WebhookTestEndpoint struct {
	ID               string
	URL              string
	SubscribedEvents []string
	IsActive         bool
}

func (r *WebhookRepository) GetWebhookEndpoint(
	ctx context.Context,
	webhookID string,
	merchantID string,
) (*WebhookTestEndpoint, error) {

	var endpoint WebhookTestEndpoint
	var events string

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			id,
			endpoint_url,
			subscribed_events,
			is_active
		FROM webhook_endpoints
		WHERE id = $1
		  AND merchant_id = $2
		LIMIT 1
		`,
		webhookID,
		merchantID,
	).Scan(
		&endpoint.ID,
		&endpoint.URL,
		&events,
		&endpoint.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("webhook endpoint not found")
		}

		return nil, fmt.Errorf(
			"failed to fetch webhook endpoint: %w",
			err,
		)
	}

	// PostgreSQL array returned as string:
	// {invoice.confirmed,invoice.expired}
	endpoint.SubscribedEvents =
		parsePostgresTextArray(events)

	return &endpoint, nil
}

func parsePostgresTextArray(value string) []string {
	if len(value) < 2 {
		return nil
	}

	value = value[1 : len(value)-1]

	if value == "" {
		return nil
	}

	var result []string
	current := ""

	for _, char := range value {
		if char == ',' {
			result = append(result, current)
			current = ""
			continue
		}

		current += string(char)
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}
