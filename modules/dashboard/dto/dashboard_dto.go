package dto

type DashboardMetricsResponse struct {
	TotalVolumeUSD          float64 `json:"total_volume_usd"`
	SuccessfulPaymentsCount int     `json:"successful_payments_count"`
	PendingInvoicesCount    int     `json:"pending_invoices_count"`
	SuccessRate             float64 `json:"success_rate"`             // Percentage, e.g., 99.4
	ActivePaymentLinksCount int     `json:"active_payment_links_count"`
}
