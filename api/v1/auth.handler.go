package v1

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// GoogleLogin godoc
// @Summary Initiate Google OAuth login
// @Description Redirects the user to the Google login page to initiate the OAuth 2.0 flow.
// @Tags auth
// @Success 302 {string} string "Redirects to Google login page"
// @Router /auth/google/login [get]
func GoogleLogin(c *fiber.Ctx) error {
	cfg := c.Locals("cfg").(*config.Config)
	url := cfg.GoogleOauthConfig.AuthCodeURL("state")
	return c.Redirect(url)
}

// GoogleCallback godoc
// @Summary Callback for Google OAuth login
// @Description Handles the callback from Google after the user has authenticated. It exchanges the authorization code for an access token, fetches the user's profile information, and then either creates a new user or logs in an existing user. Finally, it generates a JWT and returns it to the user.
// @Tags auth
// @Param code query string true "Authorization code"
// @Success 200 {object} map[string]interface{} "Returns a JWT token"
// @Router /auth/google/callback [get]
func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	db := database.DB
	cfg := c.Locals("cfg").(*config.Config)

	token, err := cfg.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return utils.InternelServerError(c, err, "Failed to exchange token")
	}

	response, err := cfg.GoogleOauthConfig.Client(context.Background(), token).Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return utils.InternelServerError(c, err, "Failed to get user info")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return utils.InternelServerError(c, err, "Failed to read user info")
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return utils.InternelServerError(c, err, "Failed to parse user info")
	}

	user, jwt, err := services.GoogleLogin(userInfo["email"].(string), userInfo["name"].(string), db, cfg)
	if err != nil {
		return utils.InternelServerError(c, err, "Failed to login with Google")
	}

	return utils.OKResponse(c, "Login successful", fiber.Map{"user": user, "token": jwt})
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with the provided full name, email, and password.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body RegisterInput true "Register Input"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	type RegisterInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return utils.BadInternalResponse(c, err, "Invalid request")
	}

	db := database.DB
	cfg := c.Locals("cfg").(*config.Config)

	user, token, err := services.Register(input.Name, input.Email, input.Password, db, cfg)
	if err != nil {
		return utils.InternelServerError(c, err, "Failed to create user")
	}

	return utils.OKCreatedResponse(c, "User registered successfully", fiber.Map{"user": user, "token": token})
}

// Login godoc
// @Summary Log in a user
// @Description Logs in a user with the provided email and password and returns a JWT token.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body LoginInput true "Login Input"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return utils.BadInternalResponse(c, err, "Invalid request")
	}

	db := database.DB
	cfg := c.Locals("cfg").(*config.Config)

	user, token, err := services.Login(input.Email, input.Password, db, cfg)
	if err != nil {
		return utils.UnauthorizedAccess(c, err, "Invalid credentials")
	}

	return utils.OKResponse(c, "Login successful", fiber.Map{"user": user, "token": token})
}

// GetProfile godoc
// @Summary Get the authenticated user's profile
// @Description Gets the profile information of the authenticated user.
// @Tags auth
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} map[string]interface{} "Profile retrieved successfully"
// @Router /auth/profile [get]
func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	db := database.DB

	user, err := services.GetProfile(userID, db)
	if err != nil {
		return utils.NotFound(c, err, "User not found")
	}

	return utils.OKResponse(c, "Profile retrieved successfully", fiber.Map{"personal": fiber.Map{"name": user.Name, "email": user.Email}})
}

// ChangePassword godoc
// @Summary Change the authenticated user's password
// @Description Changes the password of the authenticated user.
// @Tags auth
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param input body ChangePasswordInput true "Change Password Input"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Router /auth/change-password [post]
func ChangePassword(c *fiber.Ctx) error {
	type ChangePasswordInput struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	var input ChangePasswordInput

	if err := c.BodyParser(&input); err != nil {
		return utils.BadInternalResponse(c, err, "Invalid request")
	}

	userID := c.Locals("user_id").(string)
	db := database.DB

	if err := services.ChangePassword(userID, input.CurrentPassword, input.NewPassword, db); err != nil {
		return utils.InternelServerError(c, err, "Failed to change password")
	}

	return utils.OKResponse(c, "Password changed successfully", nil)
}
