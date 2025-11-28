package controllers

import (
	"gonion/src/services"

	"github.com/gofiber/fiber/v2"
)

type NpmController struct {
	Service 		*services.NpmCacheService
}

func NewNpmController(service *services.NpmCacheService) *NpmController {
	return &NpmController{
		Service: service,
	}
}

func (n* NpmController) HandleRequest(c* fiber.Ctx) error {
	cachedPath, err := n.Service.GetRessourceLocalPath(c.Params("*")) 

	if(err != nil) {
		return c.Status(400).JSON(fiber.Map{"message": "Badly formed request url"})
	}

	return c.SendFile(cachedPath)
}