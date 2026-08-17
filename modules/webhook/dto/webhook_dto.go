package dto

type RegisterWebhookEndpointRequest struct {
	EndpointURL string   `json:"endpoint_url" validate:"required,url"`
	SubscribedEvents []string `json:"subscribed_events" validate:"required"` // e.g., ["invoice.confirmed", "invoice.expired"]
}

type WebhookEventPayload struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"` // invoice.confirmed, invoice.expired
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type WebhookLogResponse struct {
	ID           string `json:"id"`
	EventType    string `json:"event_type"`
	EndpointURL  string `json:"endpoint_url"`
	ResponseCode int    `json:"response_code"`
	Status       string `json:"status"` // SUCCESS, FAILED
	CreatedAt    string `json:"created_at"`
}
