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
	UserAuthHandler 	   *handlers.UserAuthHandler
	UserService            *services.UserService
	UserRepo               *repos.UserRepo
	UserHandler            *handlers.UserHandler
	WorkspaceHandler       *handlers.WorkspaceHandler
	DefaultLimiter         ratelimiter.RateLimiter
	DB                     *pgxpool.Pool
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

	sessionStore := utilities.NewRedisSessionStore(redisClient)
	
	authLimiter, err := ratelimiter.NewRedisRateLimiter(redisClient, 1*time.Minute, 5)
	if err != nil {
		panic("failed to create rate limiter: " + err.Error())
	}
	limiter, err := ratelimiter.NewRedisRateLimiter(redisClient, 10*time.Second, 60)
	if err != nil {
		panic("failed to create rate limiter: " + err.Error())
	}
	userRepo := repos.NewUserRepo(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, sessionStore, limiter)
	
	// Initialize utilities with user service
	utilities.SetUserService(userService)
	
	userAuthHandler := handlers.NewUserAuthHandler(userService, sessionStore, authLimiter)
	
	workspaceRepo := repos.NewWorkspaceRepo(db)
	workspaceHandler := handlers.NewWorkspaceHandler(workspaceRepo, sessionStore, limiter)
	return &Container{
		AdminPanelPasswordHash: adminPanelPasswordHash,
		DB:                     db,
		AdminDashboardHandler:  handlers.NewAdminDashboardHandler(sessionStore, limiter, adminPanelPasswordHash, userService, workspaceRepo),
		AdminAuthHandler:       handlers.NewAdminAuthHandler(adminPanelPasswordHash, sessionStore, limiter),
		SessionStore:           sessionStore,
		UserService:            userService,
		UserHandler:            userHandler,
		UserRepo:               userRepo,
		UserAuthHandler:        userAuthHandler,
		DefaultLimiter:         limiter,
		WorkspaceHandler:       workspaceHandler,
	}
}
