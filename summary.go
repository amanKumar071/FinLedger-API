package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"finance_Data/internal/database"
	"finance_Data/internal/models"
)

// GetSummaryHandler returns financial summary with income, expense, and category breakdown
func GetSummaryHandler(w http.ResponseWriter, r *http.Request) {
	summary := models.SummaryResponse{
		CategoryWiseAmount: make(map[string]float64),
	}

	// Query for totals
	query := `SELECT
		SUM(IF(type='income', amount, 0)),
		SUM(IF(type='expense', amount, 0)),
		COUNT(*)
	FROM records WHERE isDeleted = 0`

	err := database.DB.QueryRow(query).Scan(&summary.Income, &summary.Expense, &summary.TotalRecord)
	if err != nil {
		http.Error(w, "Error calculating totals", http.StatusInternalServerError)
		return
	}
	summary.NetIncome = summary.Income - summary.Expense

	// Query for category totals
	rows, err := database.DB.Query("SELECT category, SUM(amount) FROM records WHERE isDeleted = 0 GROUP BY category")
	if err != nil {
		http.Error(w, "Error fetching category data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var amount float64
		if err := rows.Scan(&category, &amount); err != nil {
			http.Error(w, "Error scanning category data", http.StatusInternalServerError)
			return
		}
		summary.CategoryWiseAmount[category] = amount
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// CreateUserHandler creates a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Insert into DB
	query := "INSERT INTO users (name, role, email, status) VALUES (?, ?, ?, ?)"
	result, err := database.DB.Exec(query, user.Name, user.Role, user.Email, user.Status)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting user ID", http.StatusInternalServerError)
		return
	}
	user.ID = fmt.Sprintf("%d", id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
