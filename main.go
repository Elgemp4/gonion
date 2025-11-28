package main

import (
	"gonion/src/controllers"
	"gonion/src/repositories"
	"gonion/src/services"

	"github.com/gofiber/fiber/v2"
)

func NpmProxy() {
	npmUrl := "https://registry.npmjs.org/"
	localUrl := "http://localhost:3473/"

	npmRepository := repositories.NewNpmRessourceRepository(npmUrl)
	fsRepository := repositories.NewFsRessourceRepository("./cache/npm", npmUrl, localUrl)
	cacheRepository := repositories.NewCacheRessourceRepository(fsRepository, npmRepository)

	npmService := services.NewNpmCacheService(cacheRepository)

	npm_controller := controllers.NewNpmController(npmService)

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