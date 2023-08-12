package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kavehjamshidi/arvan-challenge/api/handler"
)

func SetupServer(
	rateLimitHandler handler.RateLimitHandler,
	uploadHandler handler.UploadHandler,
) *fiber.App {
	app := fiber.New()

	app.Post("/upload", rateLimitHandler.HandleRateLimit, uploadHandler.HandleUpload)

	return app
}
