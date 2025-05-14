package media

import (
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	service service
}

func newHandler(service service) handler {
	return handler{
		service: service,
	}
}

func (h handler) UploadHandler(c *fiber.Ctx) (err error) {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get file",
		})
	}

	upload, err := h.service.Upload(&fiber.Ctx{}, file)

	if err != nil {
		return c.Status(422).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Upload successfully",
		"data":    upload,
	})
}
func (h handler) UpdateHandler(c *fiber.Ctx) (err error) {
	req := updateRequest{}
	Id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.Update(c.Context(), req, Id)
	if errors != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errors,
		})
	}
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "data successfully updated",
		"data":    result,
	})
}
func (h handler) ReadHandler(c *fiber.Ctx) (err error) {
	Id := c.Params("id")

	result, err := h.service.Read(c.Context(), Id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "data not found"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}
func (h handler) IndexHandler(c *fiber.Ctx) (err error) {
	req := ParamIndex{}
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, err := h.service.Index(c.Context(), req)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "data not found"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) DeleteHandler(c *fiber.Ctx) (err error) {
	Id := c.Params("id")

	err = h.service.Delete(c.Context(), Id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "data successfully deleted",
	})
}
