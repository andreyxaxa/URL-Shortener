package v1

import (
	"github.com/andreyxaxa/URL-Shortener/internal/usecase"
	"github.com/andreyxaxa/URL-Shortener/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func NewLinkRoutes(apiV1Group fiber.Router, lk usecase.Link, l logger.Interface, baseURL string) {
	r := &V1{lk: lk, l: l, baseURL: baseURL}

	{
		// API
		apiV1Group.Post("/shorten", r.createShortURL)
		apiV1Group.Get("/s/:short", r.redirectToOriginalURL)
		apiV1Group.Get("/analytics/:short", r.getAnalytics)

		// Web
		apiV1Group.Get("/web", r.showUI)
	}
}
