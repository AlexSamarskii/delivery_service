package context

import (
	"github.com/AlexSamarskii/delivery_service/internal/config"
	"github.com/AlexSamarskii/delivery_service/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IDbContext interface {
	NewDBConnection() *gorm.DB
}

type IAuthMiddleware interface {
	Auth() gin.HandlerFunc
}

type AppContext struct {
	Middleware IAuthMiddleware
	DB         IDbContext
	Config     *config.Config
}

func NewAppContext(db *gorm.DB) *AppContext {
	config := config.Load()

	middlewareProvider := middleware.NewAuthMiddleware()

	dbContext := NewDbContext(db)

	return &AppContext{
		Middleware: middlewareProvider,
		DB:         dbContext,
		Config:     config,
	}
}
