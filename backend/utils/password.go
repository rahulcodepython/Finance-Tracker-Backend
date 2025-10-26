// Package utils provides utility functions for the application.
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash from the password with a cost of 14.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent.
func CheckPasswordHash(password, hash string) bool {
	// Compare the provided password with the stored hash.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// If the comparison is successful, err will be nil.
	return err == nil
}