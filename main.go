package main

import (
	"fmt"
	"log"
	"net/http"

	"finance_Data/internal/config"
	"finance_Data/internal/database"
	"finance_Data/internal/handlers"
	"finance_Data/internal/middleware"
)
func main() {
    cfg := config.NewDBConfig()
    if err := database.InitDB(cfg); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer database.CloseDB()

    mux := http.NewServeMux()

    // Chain: Auth first, then role check
    mux.Handle("/record/summary",
        middleware.AuthMiddleware(
            middleware.RequireRole("admin", "viewer","analyst")(
                http.HandlerFunc(handlers.GetSummaryHandler))))

    mux.Handle("/user/create",
        middleware.AuthMiddleware(
            middleware.RequireRole("admin")(
                http.HandlerFunc(handlers.CreateUserHandler))))

    fmt.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
