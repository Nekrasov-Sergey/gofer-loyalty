package types

type ResponseOrder struct {
	Number     string  `json:"number"`
	UploadedAt string  `json:"uploaded_at"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
}
