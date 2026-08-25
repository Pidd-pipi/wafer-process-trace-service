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

func (l Lot) IsPaused() bool { return l.Status == "hold" }

func (l Lot) IsActive() bool { return l.Status == "running" }
