// Package v1 provides handlers for version 1 of the API.
package v1

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// GoogleLogin initiates the Google OAuth 2.0 login flow.
// @Summary      Initiate Google OAuth login
// @Description  Redirects the user to Google's authentication page to begin the OAuth 2.0 flow.
// @Tags         Authentication
// @Success      302  {string}  string  "Redirects to Google login page"
// @Router       /auth/google/login [get]
func GoogleLogin(c *fiber.Ctx) error {
	// Get the Google OAuth2 configuration from the context.
	cfg := c.Locals("cfg").(*config.Config)
	// Generate the URL for the Google login page.
	url := cfg.GoogleOauthConfig.AuthCodeURL("state")
	// Redirect the user to the Google login page.
	return c.Redirect(url)
}

// GoogleCallback handles the callback from Google's OAuth 2.0 flow.
// @Summary      Callback for Google OAuth login
// @Description  Handles the callback from Google after user authentication. It exchanges the authorization code for an access token, fetches user information, and then logs in or registers the user. It returns a JWT token upon success.
// @Tags         Authentication
// @Param        code query     string true "Authorization code from Google"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /auth/google/callback [get]
func GoogleCallback(c *fiber.Ctx) error {
	// Get the authorization code from the query parameters.
	code := c.Query("code")

	// Get the database connection and configuration from the context.
	db := database.DB
	cfg := c.Locals("cfg").(*config.Config)

	// Exchange the authorization code for an access token.
	token, err := cfg.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to exchange token")
	}

	// Use the access token to get the user's profile information from Google.
	response, err := cfg.GoogleOauthConfig.Client(context.Background(), token).Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get user info")
	}

	defer response.Body.Close()

	// Read the user information from the response body.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to read user info")
	}

	// Parse the user information into a map.
	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return utils.InternalServerError(c, err, "Failed to parse user info")
	}

	// Log in or register the user with the retrieved information.
	user, jwt, err := services.GoogleLogin(userInfo["email"].(string), userInfo["name"].(string), db, cfg)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to login with Google")
	}

	// Return a 200 OK response with the user and JWT token.
	return utils.OKResponse(c, "Login successful", fiber.Map{"user": user, "token": jwt})
}

// RegisterInput defines the input structure for user registration.
type RegisterInput struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"password123"`
}

// Register handles new user registration.
// @Summary      Register a new user
// @Description  Registers a new user with a name, email, and password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        input body      RegisterInput true "Registration Input"
// @Success      201  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /auth/register [post]
func Register(c *fiber.Ctx) error {
	// Create a new instance of the RegisterInput struct.
	var input RegisterInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the database connection and configuration from the context.
	db := database.DB
	cfg := c.Locals("cfg").(*config.Config)

	// Call the Register service to create the new user.
	user, token, err := services.Register(input.Name, input.Email, input.Password, db, cfg)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create user")
	}

	// Return a 201 Created response with the new user and JWT token.
	return utils.OKCreatedResponse(c, "User registered successfully", fiber.Map{"user": user, "token": token})
}

// LoginInput defines the input structure for user login.
type LoginInput struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"password123"`
}

// Login handles user authentication.
// @Summary      Log in a user
// @Description  Authenticates a user with an email and password, returning a JWT token upon success.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        input body      LoginInput true "Login Input"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Router       /auth/login [post]
func Login(c *fiber.Ctx) error {
	// Create a new instance of the LoginInput struct.
	var input LoginInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the database connection and configuration from the context.
	db := database.DB
	cfg := c.Locals("cfg").(*config.Config)

	// Call the Login service to authenticate the user.
	user, token, err := services.Login(input.Email, input.Password, db, cfg)
	if err != nil {
		return utils.UnauthorizedAccess(c, err, "Invalid credentials")
	}

	// Return a 200 OK response with the user and JWT token.
	return utils.OKResponse(c, "Login successful", fiber.Map{"user": user, "token": token})
}

// GetProfile handles the retrieval of the authenticated user's profile.
// @Summary      Get user profile
// @Description  Retrieves the profile information of the currently authenticated user.
// @Tags         Authentication
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /auth/profile [get]
func GetProfile(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the GetProfile service to retrieve the user's profile.
	user, err := services.GetProfile(userID, db)
	if err != nil {
		return utils.NotFound(c, err, "User not found")
	}

	// Return a 200 OK response with the user's profile information.
	return utils.OKResponse(c, "Profile retrieved successfully", fiber.Map{"personal": fiber.Map{"name": user.Name, "email": user.Email}})
}

// ChangePasswordInput defines the input structure for changing a user's password.
type ChangePasswordInput struct {
	CurrentPassword string `json:"currentPassword" example:"password123"`
	NewPassword     string `json:"newPassword" example:"newpassword456"`
}

// ChangePassword handles changing the authenticated user's password.
// @Summary      Change user password
// @Description  Allows the authenticated user to change their password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        input body      ChangePasswordInput true "Change Password Input"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /auth/change-password [post]
func ChangePassword(c *fiber.Ctx) error {
	// Create a new instance of the ChangePasswordInput struct.
	var input ChangePasswordInput

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

	// Call the ChangePassword service to change the user's password.
	if err := services.ChangePassword(userID, input.CurrentPassword, input.NewPassword, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to change password")
	}

	// Return a 200 OK response.
	return utils.OKResponse(c, "Password changed successfully", nil)
}
