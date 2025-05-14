package admins

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

func (h handler) CreateHandler(c *fiber.Ctx) (err error) {
	var req = Request{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.Create(c.Context(), req)
	if errors != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errors,
		})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "new data successfully created",
		"data":    result,
	})
}

func (h handler) IndexHandler(c *fiber.Ctx) (err error) {
	var req = ParamIndex{}

	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, err := h.service.Index(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) UpdateHandler(c *fiber.Ctx) (err error) {
	var req = Request{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.Update(c.Context(), req, id)
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

func (h handler) UpdatePasswordHandler(c *fiber.Ctx) (err error) {
	var req = Request{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.UpdatePassword(c.Context(), req, id)
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

func (h handler) UpdateStatusHandler(c *fiber.Ctx) (err error) {
	var req = Request{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.UpdateStatus(c.Context(), req, id)
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
	id := c.Params("id")

	result, err := h.service.Read(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) DeleteHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	err = h.service.Delete(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "data deleted successfully",
	})
}

func (h handler) RoleIndexHandler(c *fiber.Ctx) (err error) {
	result, err := h.service.RoleIndex(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) RoleCreateHandler(c *fiber.Ctx) (err error) {
	var req = RoleRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.RoleCreate(c.Context(), req)
	if errors != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errors,
		})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "new data successfully created",
		"data":    result,
	})
}

func (h handler) RoleUpdateHandler(c *fiber.Ctx) (err error) {
	var req = RoleRequest{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.RoleUpdate(c.Context(), req, id)
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

func (h handler) RoleDeleteHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	err = h.service.RoleDelete(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "data deleted successfully",
	})
}
