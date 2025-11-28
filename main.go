package main

import (
	"gonion/controllers"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func NpmProxy() {
	npm_registry := "https://registry.npmjs.org/"
	local_registry := "http://localhost:3473/"

	client := resty.New()
	client.BaseURL = npm_registry

	npm_controller := controllers.NewNpmController(client, "./cache/npm", npm_registry, local_registry)

	app := fiber.New()

	app.Get("*", npm_controller.HandleRequest)

	app.Listen(":3473")
}

func ComposerProxy() {
	app := fiber.New()

	app.Get("/", func(c* fiber.Ctx) error {
		return c.SendString("Hello from composer proxy")
	})


	app.Listen(":3500")
}

func main() {
	go NpmProxy()
	go ComposerProxy()

	select {}
}