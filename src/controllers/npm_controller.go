package controllers

import (
	"errors"
	"gonion/src/domain"
	"gonion/src/services"
	"log/slog"
	"os"

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
	var pathError *os.PathError
	if(err != nil) {
		if(errors.As(err, &pathError)){
			slog.Error("Error while writing on the disk", "requested", c.Params("*"))
			slog.Debug(err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "Error while writing on the disk"})
		}
		if(errors.Is(err, domain.ErrRessourceCachingFailure)){
			slog.Error("Error while trying to cache the request ressource", "requested", c.Params("*"))
			return c.Status(500).JSON(fiber.Map{"message": "Error while trying to cache the request ressource"})
		}
		if(errors.Is(err, domain.ErrUnreachableUpstream)){
			slog.Error("Error from upstream while trying to fetch the request ressource", "requested", c.Params("*"))
			return c.Status(500).JSON(fiber.Map{"message": "Error from upstream while trying to fetch the request ressource"})
		}
		if(errors.Is(err, domain.ErrNotFound)){
			slog.Error("The ressource could not be found", "requested", c.Params("*"))
			return c.Status(404).JSON(fiber.Map{"message": "The ressource could not be found"})
		}
		if(errors.Is(err, domain.ErrBadRessourceName)) {
			slog.Error("Error while writing on the disk", "requested", c.Params(""))
			return c.Status(400).JSON(fiber.Map{"message": "Badly formed request url"})
		}
	}

	return c.SendFile(cachedPath)
}