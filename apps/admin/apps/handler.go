package apps

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

func (h handler) UpdateDbHandler(c *fiber.Ctx) (err error) {
	var req = Request{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.UpdateDb(c.Context(), req, id)
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
func (h handler) ReadDataHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	result, err := h.service.ReadData(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}
func (h handler) DbTestHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	result, err := h.service.DbTest(c.Context(), id)
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

func (h handler) PluginIndexHandler(c *fiber.Ctx) (err error) {
	var req = ParamIndex{}

	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, err := h.service.PluginIndex(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) PluginCreateHandler(c *fiber.Ctx) (err error) {
	var req = PluginRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.PluginCreate(c.Context(), req)
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

func (h handler) PluginUpdateHandler(c *fiber.Ctx) (err error) {
	var req = PluginRequest{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.PluginUpdate(c.Context(), req, id)
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

func (h handler) PluginDeleteHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	err = h.service.PluginDelete(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "data deleted successfully",
	})
}

func (h handler) PluginReadHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	result, err := h.service.PluginRead(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) ModulIndexHandler(c *fiber.Ctx) (err error) {
	var req = ParamIndex{}

	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, err := h.service.ModulIndex(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}
func (h handler) ModulCreateHandler(c *fiber.Ctx) (err error) {
	var req = ModulRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.ModulCreate(c.Context(), req)
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
func (h handler) ModulUpdateHandler(c *fiber.Ctx) (err error) {
	var req = ModulRequest{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.ModulUpdate(c.Context(), req, id)
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
func (h handler) ModulUpdateStatusHandler(c *fiber.Ctx) (err error) {
	var req = ModulRequest{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.ModulUpdateStatus(c.Context(), req, id)
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
func (h handler) ModulDeleteHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	err = h.service.ModulDelete(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "data deleted successfully",
	})
}
func (h handler) ModulReadHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	result, err := h.service.ModulRead(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func (h handler) FeatureCreateHandler(c *fiber.Ctx) (err error) {
	var req = FeatureRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.FeatureCreate(c.Context(), req)
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
func (h handler) FeatureUpdateHandler(c *fiber.Ctx) (err error) {
	var req = FeatureRequest{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.FeatureUpdate(c.Context(), req, id)
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
func (h handler) FeatureUpdateStatusHandler(c *fiber.Ctx) (err error) {
	var req = FeatureRequest{}
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result, errors, err := h.service.FeatureUpdateStatus(c.Context(), req, id)
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
func (h handler) FeatureDeleteHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	err = h.service.FeatureDelete(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "data deleted successfully",
	})
}
func (h handler) FeatureReadHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	result, err := h.service.FeatureRead(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}
func (h handler) FeatureBulkDeleteHandler(c *fiber.Ctx) (err error) {
	var req = Request{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	success, failed, err := h.service.FeatureBulkDelete(c.Context(), req)
	if err != nil {
		return c.Status(422).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "data deleted successfully",
		"data": fiber.Map{
			"success": success,
			"failed":  failed,
		},
	})
}
