package types

type OrderResponse struct {
	Number     string      `json:"number"`
	UploadedAt string      `json:"uploaded_at"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual,omitempty"`
}
