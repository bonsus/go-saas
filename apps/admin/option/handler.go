package option

import (
	"github.com/gofiber/fiber/v2"
)

// var svc service

type handler struct {
	service service
}

func newHandler(service service) handler {
	return handler{
		service: service,
	}
}

func (h handler) updateHandler(c *fiber.Ctx) (err error) {
	var req = Request{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}
	data, errors, err := h.service.Save(c.Context(), req)
	if err != nil {
		if errors != nil {
			return c.Status(422).JSON(fiber.Map{
				"message": err.Error(),
				"errors":  errors,
			})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}
	return c.Status(200).JSON(fiber.Map{
		"data": data.Value,
	})
}

func (h handler) getHandler(c *fiber.Ctx) (err error) {
	name := c.Params("name")
	data, err := h.service.Get(name)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Option not found"})
	}
	return c.Status(200).JSON(fiber.Map{
		"data": data.Value,
	})
}
