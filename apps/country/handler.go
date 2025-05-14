package country

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

func (h handler) CountryHandler(c *fiber.Ctx) (err error) {
	var req = Params{}

	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, err := h.service.Countries(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) ProvinceHandler(c *fiber.Ctx) (err error) {
	var req = Params{}

	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, err := h.service.Provinces(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) CityHandler(c *fiber.Ctx) (err error) {
	var req = Params{}

	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	if req.Province == "" {
		return c.Status(422).JSON(fiber.Map{"errors": "province is required"})
	}

	result, err := h.service.Cities(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}
func (h handler) DistrictHandler(c *fiber.Ctx) (err error) {
	var req = Params{}

	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	if req.Province == "" {
		return c.Status(422).JSON(fiber.Map{"errors": "province is required"})
	}
	if req.City == "" {
		return c.Status(422).JSON(fiber.Map{"errors": "city is required"})
	}

	result, err := h.service.Districts(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error", "ss": err})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) ZipHandler(c *fiber.Ctx) (err error) {
	var req = Params{}

	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	if req.Province == "" {
		return c.Status(422).JSON(fiber.Map{"error": "province is required"})
	}
	if req.City == "" {
		return c.Status(422).JSON(fiber.Map{"error": "city is required"})
	}
	if req.District == "" {
		return c.Status(422).JSON(fiber.Map{"error": "district is required"})
	}

	result, err := h.service.Zips(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) SearchHandler(c *fiber.Ctx) (err error) {
	var req = Params{}

	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, err := h.service.Search(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}
