package models

// User represents a user in the system
type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// SummaryResponse contains financial summary data
type SummaryResponse struct {
	Expense            float64            `json:"expense"`
	Income             float64            `json:"income"`
	NetIncome          float64            `json:"netIncome"`
	TotalRecord        int                `json:"totalRecord"`
	CategoryWiseAmount map[string]float64 `json:"categoryWiseAmount"`
}
