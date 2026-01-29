package v1

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var (
	//go:embed web/index.html
	webFile embed.FS
)

func (r *V1) showUI(ctx *fiber.Ctx) error {
	file, err := webFile.ReadFile("web/index.html")
	if err != nil {
		r.l.Error(err, "restapi - v1 - showUI")

		return errorResponse(ctx, http.StatusInternalServerError, "problems with load UI")
	}

	ctx.Set("Content-Type", "text/html")
	return ctx.Send(file)
}
