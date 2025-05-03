package di

import (
	"backend/internal/handlers"
	"backend/internal/repos"
	"backend/internal/services"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Container struct {
	AdminPanelPasswordHash []byte
	AdminDashboardHandler  *handlers.AdminDashboardHandler
	AdminAuthHandler       *handlers.AdminAuthHandler
	UserService            *services.UserService
	UserRepo               *repos.UserRepo
	UserHandler            *handlers.UserHandler
	DB                     *pgxpool.Pool
	Limiter                ratelimiter.RateLimiter
	SessionStore           utilities.SessionStore
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

	limiter, err := ratelimiter.NewRedisRateLimiter(redisClient, 10*time.Second, 10)
	sessionStore := utilities.NewRedisSessionStore(redisClient)
	if err != nil {
		panic("failed to create rate limiter: " + err.Error())
	}

	userRepo := repos.NewUserRepo(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, sessionStore, limiter)
	return &Container{
		AdminPanelPasswordHash: adminPanelPasswordHash,
		DB:                     db,
		AdminDashboardHandler:  handlers.NewAdminDashboardHandler(sessionStore, limiter, adminPanelPasswordHash, userService),
		AdminAuthHandler:       handlers.NewAdminAuthHandler(adminPanelPasswordHash, sessionStore, limiter),
		Limiter:                limiter,
		SessionStore:           sessionStore,
		UserService:            userService,
		UserHandler:            userHandler,
		UserRepo:               userRepo,
	}
}
