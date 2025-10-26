// Package v1 provides handlers for version 1 of the API.
package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateCategoryInput defines the input structure for creating a new transaction category.
type CreateCategoryInput struct {
	Name string `json:"name" example:"Groceries"`
	Type string `json:"type" example:"Expense"`
}

// CreateCategory handles the creation of a new transaction category.
// @Summary      Create a new transaction category
// @Description  Creates a new transaction category for the authenticated user.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        input body      CreateCategoryInput true "Create Category Input"
// @Success      201  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /categories/create [post]
func CreateCategory(c *fiber.Ctx) error {
	// Create a new instance of the CreateCategoryInput struct.
	var input CreateCategoryInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the CreateCategory service to create the new category.
	category, err := services.CreateCategory(input.Name, models.TransactionType(input.Type), userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create category")
	}

	// Return a 201 Created response with the new category.
	return utils.OKCreatedResponse(c, "Category created successfully", category)
}

// GetCategories handles the retrieval of all transaction categories.
// @Summary      Get all transaction categories
// @Description  Retrieves all transaction categories available in the system.
// @Tags         Categories
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /categories [get]
func GetCategories(c *fiber.Ctx) error {
	// Get the database connection.
	db := database.DB
	// Call the GetCategories service to retrieve all categories.
	categories, err := services.GetCategories(db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to retrieve categories")
	}

	// Return a 200 OK response with the retrieved categories.
	return utils.OKResponse(c, "Categories retrieved successfully", categories)
}

// UpdateCategoryInput defines the input structure for updating a transaction category.
type UpdateCategoryInput struct {
	Name string `json:"name" example:"Food"`
	Type string `json:"type" example:"Expense"`
}

// UpdateCategory handles the update of a specific transaction category.
// @Summary      Update a transaction category
// @Description  Updates the details of a specific transaction category for the authenticated user.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id   path      string            true "Category ID"
// @Param        input body      UpdateCategoryInput true "Update Category Input"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /categories/update/{id} [patch]
func UpdateCategory(c *fiber.Ctx) error {
	// Create a new instance of the UpdateCategoryInput struct.
	var input UpdateCategoryInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the category ID from the URL parameters.
	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid category ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the UpdateCategory service to update the category.
	category, err := services.UpdateCategory(categoryID, input.Name, models.TransactionType(input.Type), userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to update category")
	}

	// Return a 200 OK response with the updated category.
	return utils.OKResponse(c, "Category updated successfully", category)
}

// DeleteCategory handles the deletion of a specific transaction category.
// @Summary      Delete a transaction category
// @Description  Deletes a specific transaction category for the authenticated user.
// @Tags         Categories
// @Produce      json
// @Param        id   path      string  true "Category ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /categories/delete/{id} [delete]
func DeleteCategory(c *fiber.Ctx) error {
	// Get the category ID from the URL parameters.
	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid category ID")
	}

	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the DeleteCategory service to delete the category.
	if err := services.DeleteCategory(categoryID, userID, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to delete category")
	}

	// Return a 200 OK response.
	return utils.OKResponse(c, "Category deleted successfully", nil)
}
