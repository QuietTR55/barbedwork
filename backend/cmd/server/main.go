package main

import (
	"backend/internal/di"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	fmt.Println("Starting server...")

	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Connect to the database
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Dependency injection container
	container := di.NewContainer(db)

	// Serve static files from the "./uploads" directory
	uploadsDir := "./uploads/"
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	// Setup routes
	SetupRoutes(mux, container)

	// Conditionally apply CORS middleware
	devMode := os.Getenv("DEV") != ""
	var wrappedMux http.Handler
	if devMode {
		fmt.Println("Applying CORS middleware")
		wrappedMux = corsMiddleware(mux)
	} else {
		wrappedMux = mux
	}

	// Start the server
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, wrappedMux))
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" { // Replace with your frontend URL
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
