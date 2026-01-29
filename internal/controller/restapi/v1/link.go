package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/andreyxaxa/URL-Shortener/internal/controller/restapi/v1/request"
	"github.com/andreyxaxa/URL-Shortener/internal/controller/restapi/v1/response"
	"github.com/andreyxaxa/URL-Shortener/internal/controller/restapi/v1/validate"
	"github.com/andreyxaxa/URL-Shortener/pkg/types/errs"
	"github.com/gofiber/fiber/v2"
)

type analyticsHandler func(ctx *fiber.Ctx) error

func (r *V1) createShortURL(ctx *fiber.Ctx) error {
	var body request.CreateShortURLRequest

	err := ctx.BodyParser(&body)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid url")
	}

	if !validate.IsValidURL(body.URL) {
		return errorResponse(ctx, http.StatusBadRequest, "invalid url")
	}

	if body.CustomAlias != "" {
		if !validate.IsValidAlias(body.CustomAlias) {
			return errorResponse(ctx, http.StatusBadRequest, "invalid alias: use only letters, numbers, dash, underscope")
		}

		if !validate.IsValidAliasLength(body.CustomAlias) {
			return errorResponse(ctx, http.StatusBadRequest, "invalid alias length: must be 3-50 chars")
		}
	}

	shortCode, err := r.lk.CreateShortURL(ctx.UserContext(), body.URL, body.CustomAlias)
	if err != nil {
		if errors.Is(err, errs.ErrAliasAlreadyTaken) {
			return errorResponse(ctx, http.StatusBadRequest, "alias already taken")
		}
		r.l.Error(err, "restapi - v1 - createShortURL")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	resp := response.CreateShortURLResponse{
		URL:      body.URL,
		ShortURL: fmt.Sprintf("%s/s/%s", r.baseURL, shortCode),
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func (r *V1) redirectToOriginalURL(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("short")

	originalURL, err := r.lk.GetOriginalURLByShortCode(ctx.UserContext(), shortCode)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "couldnt find original URL")
		}
		r.l.Error(err, "restapi - v1 - redirectToOriginalURL")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	err = r.lk.TrackClick(ctx.UserContext(), shortCode, ctx.IP(), ctx.Get("User-Agent"))
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "couldnt find original URL")
		}
		r.l.Error(err, "restapi - v1 - redirectToOriginalURL")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	return ctx.Redirect(originalURL, http.StatusMovedPermanently)
}

func (r *V1) getAnalytics(ctx *fiber.Ctx) error {
	groupBy := ctx.Query("group-by")

	strats := map[string]analyticsHandler{
		"":        r.getFullAnalytics,
		"day":     r.getAnalyticsByDate,
		"month":   r.getAnalyticsByDate,
		"device":  r.getAnalyticsByDevice,
		"browser": r.getAnalyticsByBrowser,
	}

	handler, ok := strats[groupBy]
	if !ok {
		return errorResponse(ctx, http.StatusBadRequest, "invalid group-by")
	}

	return handler(ctx)
}

func (r *V1) getFullAnalytics(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("short")

	err := r.lk.ExistsByShortCode(ctx.UserContext(), shortCode)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusBadRequest, "couldnt find original URL")
		}
		r.l.Error(err, "restapi - v1 - getFullAnalytics")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	fullAnalytics, err := r.lk.GetAnalytics(ctx.UserContext(), shortCode)
	if err != nil {
		r.l.Error(err, "restapi - v1 - getAnalytics")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	analytics := response.Analytics{
		TotalClicks:     fullAnalytics.TotalClicks,
		ClicksByBrowser: fullAnalytics.ClicksByBrowser,
		ClicksByDevice:  fullAnalytics.ClicksByDevice,
		RecentClicks:    make([]response.ClickByDate, 0),
	}

	resp := response.GetAnalyticsResponse{
		Analytics: analytics,
	}

	for _, a := range fullAnalytics.RecentClicks {
		resp.Analytics.RecentClicks = append(resp.Analytics.RecentClicks, response.ClickByDate{
			Date:   a.Date.Format("2006-01-02"),
			Clicks: a.Clicks,
		})
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func (r *V1) getAnalyticsByDate(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("short")
	groupBy := ctx.Query("group-by")

	err := r.lk.ExistsByShortCode(ctx.UserContext(), shortCode)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "couldnt find original URL")
		}
		r.l.Error(err, "restapi - v1 - getAnalyticsByDate")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	recentClicks, err := r.lk.GetRecentClicks(ctx.UserContext(), shortCode, groupBy)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidInterval) {
			return errorResponse(ctx, http.StatusBadRequest, "invalid interval: must be \"day\" or \"month\"")
		}
		r.l.Error(err, "restapi - v1 - getAnalyticsByDate")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	format := "2006-01-02"
	if groupBy == "month" {
		format = "2006-01"
	}

	analytics := response.AnalyticsByDate{
		RecentClicks: make([]response.ClickByDate, 0),
	}
	resp := response.GetAnalyticsByDateResponse{
		Analytics: analytics,
	}

	for _, a := range recentClicks {
		resp.Analytics.RecentClicks = append(resp.Analytics.RecentClicks, response.ClickByDate{
			Date:   a.Date.Format(format),
			Clicks: a.Clicks,
		})
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func (r *V1) getAnalyticsByBrowser(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("short")

	err := r.lk.ExistsByShortCode(ctx.UserContext(), shortCode)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "couldnt find original URL")
		}
		r.l.Error(err, "restapi - v1 - getAnalyticsByBrowser")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	clicksByBrowser, err := r.lk.GetClicksByBrowser(ctx.UserContext(), shortCode)
	if err != nil {
		r.l.Error(err, "restapi - v1 - getAnalyticsByBrowser")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	resp := response.GetAnalyticsByBrowserResponse{
		Analytics: response.AnalyticsByBrowser{ClicksByBrowser: clicksByBrowser},
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func (r *V1) getAnalyticsByDevice(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("short")

	err := r.lk.ExistsByShortCode(ctx.UserContext(), shortCode)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "couldnt find original URL")
		}
		r.l.Error(err, "restapi - v1 - getAnalyticsByDevice")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	clicksByDevice, err := r.lk.GetClicksByDevice(ctx.UserContext(), shortCode)
	if err != nil {
		r.l.Error(err, "restapi - v1 - getAnalyticsByDevice")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	resp := response.GetAnalyticsByDeviceResponse{
		Analytics: response.AnalyticsByDevice{ClicksByDevice: clicksByDevice},
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}
