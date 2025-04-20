package main

import (
	"backend/internal/di"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	fmt.Println("Starting server...")
	router := gin.Default()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	container := di.NewContainer(db)

	SetupRoutes(router, container)

	router.Run(":" + port)
}
