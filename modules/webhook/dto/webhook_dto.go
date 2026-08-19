package dto

// ============================================================
// REGISTER WEBHOOK
// ============================================================

type RegisterWebhookEndpointRequest struct {
	EndpointURL string `json:"endpoint_url" validate:"required,url"`

	SubscribedEvents []string `json:"subscribed_events" validate:"required,min=1"`
}

// ============================================================
// REGISTER WEBHOOK RESPONSE
// ============================================================
//
// The secret should ONLY be returned when the webhook is
// created/registered.
//
// Do not expose the secret through normal GET/list endpoints.
//

type RegisterWebhookEndpointResponse struct {
	ID string `json:"id"`

	EndpointURL string `json:"endpoint_url"`

	SubscribedEvents []string `json:"subscribed_events"`

	Secret string `json:"secret"`

	Active bool `json:"active"`

	CreatedAt string `json:"created_at"`
}

// ============================================================
// WEBHOOK EVENT
// ============================================================

type WebhookEventPayload struct {
	EventID string `json:"event_id"`

	EventType string `json:"event_type"`

	Timestamp int64 `json:"timestamp"`

	Data interface{} `json:"data"`
}

// ============================================================
// WEBHOOK LOG
// ============================================================

type WebhookLogResponse struct {
	ID string `json:"id"`

	EventID string `json:"event_id"`

	EventType string `json:"event_type"`

	EndpointURL string `json:"endpoint_url"`

	ResponseCode int `json:"response_code"`

	Status string `json:"status"`

	CreatedAt string `json:"created_at"`
}

// ============================================================
// WEBHOOK ENDPOINT RESPONSE
// ============================================================
//
// Used for GET /webhooks.
// Notice that Secret is intentionally NOT included.
//

type WebhookEndpointResponse struct {
	ID string `json:"id"`

	EndpointURL string `json:"endpoint_url"`

	SubscribedEvents []string `json:"subscribed_events"`

	Active bool `json:"active"`

	CreatedAt string `json:"created_at"`
}
