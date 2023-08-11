package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
)

func main() {
	app := fiber.New()

	address := viper.GetString("SERVER_ADDRESS")
	if err := app.Listen(address); err != nil {
		log.Panicf("could not start server: %v\n", err)
	}

	app.Post("/upload")
}
