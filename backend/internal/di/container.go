package di

import (
	"backend/internal/handlers"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Container struct {
	AdminPanelPasswordHash []byte
	AdminDashboardHandler  *handlers.AdminDashboardHandler
	AdminAuthHandler       *handlers.AdminAuthHandler
	DB                     *pgxpool.Pool
	RedisClient            *redis.Client
}

func NewContainer(db *pgxpool.Pool) *Container {
	panelPassword := os.Getenv("ADMIN_PANEL_PASSWORD")
	if panelPassword == "" {
		panic("please set an admin panel password for security reasons")
	}
	adminPanelPasswordHash, err := bcrypt.GenerateFromPassword([]byte(panelPassword), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to generate password hash")
	}

	// Redis connection
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &Container{
		AdminPanelPasswordHash: adminPanelPasswordHash,
		DB:                     db,
		AdminDashboardHandler:  handlers.NewAdminDashboardHandler(redisClient, adminPanelPasswordHash),
		AdminAuthHandler:       handlers.NewAdminAuthHandler(adminPanelPasswordHash, redisClient),
		RedisClient:            redisClient,
	}
}
