package auth

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

func (h handler) LoginHandler(c *fiber.Ctx) (err error) {
	var req = Request{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	data, errorsMap, err := h.service.Login(c.Context(), req)
	if err != nil {
		if errorsMap != nil {
			return c.Status(422).JSON(fiber.Map{
				"error":  err.Error(),
				"errors": errorsMap,
			})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Login Successfully",
		"data":    data,
	})
}

func (h handler) RegisterHandler(c *fiber.Ctx) (err error) {
	var req = RegisterRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	_, errorsMap, err := h.service.Register(req)
	if err != nil {
		if errorsMap != nil {
			return c.Status(422).JSON(fiber.Map{"errors": errorsMap})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "new user successfully created",
	})
}

func (h handler) meHandler(c *fiber.Ctx) (err error) {
	userId, ok := c.Locals("admin_id").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	data, err := h.service.Me(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}

func (h handler) updateHandler(c *fiber.Ctx) (err error) {
	userId, ok := c.Locals("admin_id").(string)
	req := RegisterRequest{}
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}
	data, errors, err := h.service.Update(req, userId)
	if errors != nil {
		return c.Status(422).JSON(fiber.Map{"errors": errors, "s": req})
	}
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}

func (h handler) updatePasswordHandler(c *fiber.Ctx) (err error) {
	userId, ok := c.Locals("admin_id").(string)
	req := RegisterRequest{}
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}
	data, errors, err := h.service.UpdatePassword(req, userId)
	if errors != nil {
		return c.Status(422).JSON(fiber.Map{"errors": errors})
	}
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}
