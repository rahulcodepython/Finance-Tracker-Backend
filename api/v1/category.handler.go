package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
)

// CreateCategory godoc
// @Summary Create a new transaction category
// @Description Creates a new transaction category.
// @Tags categories
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param input body CreateCategoryInput true "Create Category Input"
// @Success 201 {object} map[string]interface{} "Category created successfully"
// @Router /categories/create [post]
func CreateCategory(c *fiber.Ctx) error {
	type CreateCategoryInput struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	var input CreateCategoryInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}

	db := database.DB

	category, err := services.CreateCategory(input.Name, models.TransactionType(input.Type), db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to create category", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Category created successfully", "data": category})
}

// GetCategories godoc
// @Summary Get all transaction categories
// @Description Gets all transaction categories.
// @Tags categories
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} map[string]interface{} "Categories retrieved successfully"
// @Router /categories [get]
func GetCategories(c *fiber.Ctx) error {
	db := database.DB
	categories, err := services.GetCategories(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to get categories", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Categories retrieved successfully", "data": categories})
}

// UpdateCategory godoc
// @Summary Update a transaction category
// @Description Updates a transaction category.
// @Tags categories
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Param input body UpdateCategoryInput true "Update Category Input"
// @Success 200 {object} map[string]interface{} "Category updated successfully"
// @Router /categories/update/{id} [patch]
func UpdateCategory(c *fiber.Ctx) error {
	type UpdateCategoryInput struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	var input UpdateCategoryInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}

	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid category ID", "error": err.Error()})
	}

	db := database.DB

	category, err := services.UpdateCategory(categoryID, input.Name, models.TransactionType(input.Type), db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to update category", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Category updated successfully", "data": category})
}

// DeleteCategory godoc
// @Summary Delete a transaction category
// @Description Deletes a transaction category.
// @Tags categories
// @Security ApiKeyAuth
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{} "Category deleted successfully"
// @Router /categories/delete/{id} [delete]
func DeleteCategory(c *fiber.Ctx) error {
	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid category ID", "error": err.Error()})
	}

	db := database.DB

	if err := services.DeleteCategory(categoryID, db); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to delete category", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Category deleted successfully"})
}
