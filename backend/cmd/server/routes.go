package main

import (
	"backend/internal/di"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router gin.IRouter, container *di.Container) {
	router.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error without exposing stack traces
				c.Error(err.(error))
				c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
			}
		}()
		c.Next()
	})
	router.Use(gin.Logger())

	adminDashboard := router.Group("/admin")
	{
		adminDashboard.GET("/dashboard", container.AdminDashboardHandler.GetDashboardData)
	}
}
