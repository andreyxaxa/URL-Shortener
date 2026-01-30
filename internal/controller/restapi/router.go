package restapi

import (
	"github.com/andreyxaxa/URL-Shortener/config"
	_ "github.com/andreyxaxa/URL-Shortener/docs" // Swagger docs.
	v1 "github.com/andreyxaxa/URL-Shortener/internal/controller/restapi/v1"
	"github.com/andreyxaxa/URL-Shortener/internal/usecase"
	"github.com/andreyxaxa/URL-Shortener/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title URL Shortener
// @version 1.0
// @host localhost:8080
// @BasePath /v1
func NewRouter(app *fiber.App, cfg *config.Config, lk usecase.Link, l logger.Interface, baseURL string) {
	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewLinkRoutes(apiV1Group, lk, l, baseURL+"/v1")
	}
}
