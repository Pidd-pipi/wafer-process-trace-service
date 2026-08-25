package domain

type Lot struct {
	ID         string `json:"id"`
	WaferCount int    `json:"wafer_count"`
	Step       string `json:"step"`
	Equipment  string `json:"equipment"`
	Status     string `json:"status"`
	StartedAt  string `json:"started_at"`
	UpdatedAt  string `json:"updated_at"`
}
