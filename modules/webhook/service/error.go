package service

import "errors"

var (
	ErrInvalidWebhook      = errors.New("invalid webhook")
	ErrWebhookNotFound     = errors.New("webhook not found")
	ErrInvalidWebhookURL   = errors.New("invalid webhook URL")
	ErrUnauthorizedWebhook = errors.New("unauthorized webhook")
	ErrWebhookInactive     = errors.New("webhook is inactive")
)
