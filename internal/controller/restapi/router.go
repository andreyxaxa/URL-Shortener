package restapi

import (
	v1 "github.com/andreyxaxa/URL-Shortener/internal/controller/restapi/v1"
	"github.com/andreyxaxa/URL-Shortener/internal/usecase"
	"github.com/andreyxaxa/URL-Shortener/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, lk usecase.Link, l logger.Interface, baseURL string) {
	apiV1Group := app.Group("/v1")
	{
		v1.NewLinkRoutes(apiV1Group, lk, l, baseURL+"/v1")
	}
}
