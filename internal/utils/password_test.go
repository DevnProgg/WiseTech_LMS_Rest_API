package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "TestPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("Hashed password is empty")
	}

	// Verify that the hashed password can be checked successfully
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Errorf("CompareHashAndPassword failed for the newly hashed password: %v", err)
	}
}

func TestCheckPassword(t *testing.T) {
	password := "TestPassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		hashedPass     string
		plainPass      string
		expectError    bool
		expectedBcrypt bool // true if bcrypt.CompareHashAndPassword is expected to return nil
	}{
		{
			name:           "Correct password",
			hashedPass:     string(hashedPassword),
			plainPass:      password,
			expectError:    false,
			expectedBcrypt: true,
		},
		{
			name:           "Incorrect password",
			hashedPass:     string(hashedPassword),
			plainPass:      "WrongPassword123",
			expectError:    true,
			expectedBcrypt: false,
		},
		{
			name:           "Empty plain password",
			hashedPass:     string(hashedPassword),
			plainPass:      "",
			expectError:    true,
			expectedBcrypt: false,
		},
		{
			name:           "Empty hashed password",
			hashedPass:     "",
			plainPass:      password,
			expectError:    true,
			expectedBcrypt: false,
		},
		{
			name:           "Invalid hashed password format",
			hashedPass:     "invalidhash",
			plainPass:      password,
			expectError:    true,
			expectedBcrypt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(tt.hashedPass, tt.plainPass)
			if (err != nil) != tt.expectError {
				t.Errorf("CheckPassword() error = %v, expectError %v", err, tt.expectError)
			}
			// Further verify the error type if it's expected to be a bcrypt.ErrMismatchedHashAndPassword
			if tt.expectError && !tt.expectedBcrypt && err == nil {
				t.Errorf("CheckPassword() expected bcrypt error but got none")
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
		expectedMsg string
	}{
		{
			name:        "Valid password",
			password:    "StrongPass123",
			expectError: false,
		},
		{
			name:        "Too short",
			password:    "Short1",
			expectError: true,
			expectedMsg: "password must be at least 8 characters long",
		},
		{
			name:        "No uppercase",
			password:    "nouppercase123",
			expectError: true,
			expectedMsg: "password must contain at least one uppercase letter",
		},
		{
			name:        "No number",
			password:    "NoNumberTest",
			expectError: true,
			expectedMsg: "password must contain at least one number",
		},
		{
			name:        "No uppercase and too short", // Expecting the first error to be returned
			password:    "short1",
			expectError: true,
			expectedMsg: "password must be at least 8 characters long",
		},
		{
			name:        "Empty password",
			password:    "",
			expectError: true,
			expectedMsg: "password must be at least 8 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidatePassword() error = %v, expectError %v", err, tt.expectError)
			}
			if tt.expectError && err != nil && err.Error() != tt.expectedMsg {
				t.Errorf("ValidatePassword() got error message = %q, want %q", err.Error(), tt.expectedMsg)
			}
		})
	}
}
