package di

import (
	"backend/internal/handlers"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Container struct {
	AdminPanelPasswordHash []byte
	AdminDashboardHandler  *handlers.AdminDashboardHandler
	DB                     *pgxpool.Pool
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
	return &Container{
		AdminPanelPasswordHash: adminPanelPasswordHash,
		DB:                     db,
		AdminDashboardHandler:  handlers.NewAdminDashboardHandler(),
	}
}
