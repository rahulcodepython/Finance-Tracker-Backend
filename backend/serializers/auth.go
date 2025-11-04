package serializers

// RegisterInput defines the input structure for user registration.
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginInput defines the input structure for user login.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ChangePasswordInput defines the input structure for changing a user's password.
type ChangePasswordInput struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
