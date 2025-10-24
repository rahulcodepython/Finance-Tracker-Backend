package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
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
		return utils.BadResponse(c, err, "Invalid request")
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	db := database.DB

	category, err := services.CreateCategory(input.Name, models.TransactionType(input.Type), userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create category")
	}

	return utils.OKCreatedResponse(c, "Category created successfully", category)
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
		return utils.InternalServerError(c, err, "Failed to get categories")
	}

	return utils.OKResponse(c, "Categories retrieved successfully", categories)
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
		return utils.BadResponse(c, err, "Invalid request")
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid category ID")
	}

	db := database.DB

	category, err := services.UpdateCategory(categoryID, input.Name, models.TransactionType(input.Type), userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to update category")
	}

	return utils.OKResponse(c, "Category updated successfully", category)
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
		return utils.BadResponse(c, err, "Invalid category ID")
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	db := database.DB

	if err := services.DeleteCategory(categoryID, userID, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to delete category")
	}

	return utils.OKResponse(c, "Category deleted successfully", nil)
}
