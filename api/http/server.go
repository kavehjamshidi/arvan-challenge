package http

import (
	"github.com/gofiber/fiber/v2"
	handler2 "github.com/kavehjamshidi/arvan-challenge/api/http/handler"
	"github.com/spf13/viper"
	"log"
)

func SetupServer(
	rateLimitHandler handler2.RateLimitHandler,
	uploadHandler handler2.UploadHandler,
) {
	app := fiber.New()

	app.Post("/upload", rateLimitHandler.HandleRateLimit, uploadHandler.HandleUpload)

	address := viper.GetString("SERVER_ADDRESS")
	if err := app.Listen(address); err != nil {
		log.Panicf("could not initialize server: %v\n", err)
	}
}
