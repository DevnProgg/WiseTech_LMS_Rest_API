package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "supersecretkey"
const testUserID = 123

func TestGenerateAccessToken(t *testing.T) {
	tokenString, err := GenerateAccessToken(testUserID, testSecret)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}
	if tokenString == "" {
		t.Error("Generated access token is empty")
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse access token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		t.Fatal("Access token is invalid or claims are not of type *Claims")
	}

	if claims.UserID != testUserID {
		t.Errorf("Expected UserID %d, got %d", testUserID, claims.UserID)
	}

	// Check expiration within a reasonable delta
	expectedExpiry := time.Now().Add(AccessTokenDuration)
	if !claims.ExpiresAt.After(time.Now()) || claims.ExpiresAt.After(expectedExpiry.Add(time.Second)) {
		t.Errorf("Access token expiry is not within expected range. Expected around %v, got %v", expectedExpiry, claims.ExpiresAt.Time)
	}
	if !claims.IssuedAt.Before(time.Now().Add(time.Second)) {
		t.Errorf("Access token issued at time is not correct. Expected around %v, got %v", time.Now(), claims.IssuedAt.Time)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tokenString, err := GenerateRefreshToken(testUserID, testSecret)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if tokenString == "" {
		t.Error("Generated refresh token is empty")
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse refresh token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		t.Fatal("Refresh token is invalid or claims are not of type *Claims")
	}

	if claims.UserID != testUserID {
		t.Errorf("Expected UserID %d, got %d", testUserID, claims.UserID)
	}

	// Check expiration within a reasonable delta
	expectedExpiry := time.Now().Add(RefreshTokenDuration)
	if !claims.ExpiresAt.After(time.Now()) || claims.ExpiresAt.After(expectedExpiry.Add(time.Second)) {
		t.Errorf("Refresh token expiry is not within expected range. Expected around %v, got %v", expectedExpiry, claims.ExpiresAt.Time)
	}
	if !claims.IssuedAt.Before(time.Now().Add(time.Second)) {
		t.Errorf("Refresh token issued at time is not correct. Expected around %v, got %v", time.Now(), claims.IssuedAt.Time)
	}
}

func TestGenerateTokenPair(t *testing.T) {
	tokenPair, err := GenerateTokenPair(testUserID, testSecret)
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}
	if tokenPair == nil {
		t.Fatal("Generated token pair is nil")
	}
	if tokenPair.AccessToken == "" {
		t.Error("Access token in pair is empty")
	}
	if tokenPair.RefreshToken == "" {
		t.Error("Refresh token in pair is empty")
	}

	// Validate Access Token from the pair
	accessToken, err := jwt.ParseWithClaims(tokenPair.AccessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse access token from pair: %v", err)
	}
	accessClaims, ok := accessToken.Claims.(*Claims)
	if !ok || !accessToken.Valid {
		t.Fatal("Access token from pair is invalid or claims are not of type *Claims")
	}
	if accessClaims.UserID != testUserID {
		t.Errorf("Expected UserID %d for access token, got %d", testUserID, accessClaims.UserID)
	}
	expectedAccessExpiry := time.Now().Add(AccessTokenDuration)
	if !accessClaims.ExpiresAt.After(time.Now()) || accessClaims.ExpiresAt.After(expectedAccessExpiry.Add(time.Second)) {
		t.Errorf("Access token expiry from pair is not within expected range. Expected around %v, got %v", expectedAccessExpiry, accessClaims.ExpiresAt.Time)
	}

	// Validate Refresh Token from the pair
	refreshToken, err := jwt.ParseWithClaims(tokenPair.RefreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse refresh token from pair: %v", err)
	}
	refreshClaims, ok := refreshToken.Claims.(*Claims)
	if !ok || !refreshToken.Valid {
		t.Fatal("Refresh token from pair is invalid or claims are not of type *Claims")
	}
	if refreshClaims.UserID != testUserID {
		t.Errorf("Expected UserID %d for refresh token, got %d", testUserID, refreshClaims.UserID)
	}
	expectedRefreshExpiry := time.Now().Add(RefreshTokenDuration)
	if !refreshClaims.ExpiresAt.After(time.Now()) || refreshClaims.ExpiresAt.After(expectedRefreshExpiry.Add(time.Second)) {
		t.Errorf("Refresh token expiry from pair is not within expected range. Expected around %v, got %v", expectedRefreshExpiry, refreshClaims.ExpiresAt.Time)
	}
}

func TestGenerateTokenInvalidSecret(t *testing.T) {
	expectedErr := "secret key cannot be empty"

	_, err := GenerateAccessToken(testUserID, "")
	if err == nil || err.Error() != expectedErr {
		t.Errorf("GenerateAccessToken with empty secret key: expected error %q, got %v", expectedErr, err)
	}

	_, err = GenerateRefreshToken(testUserID, "")
	if err == nil || err.Error() != expectedErr {
		t.Errorf("GenerateRefreshToken with empty secret key: expected error %q, got %v", expectedErr, err)
	}

	_, err = GenerateTokenPair(testUserID, "")
	if err == nil || err.Error() != expectedErr {
		t.Errorf("GenerateTokenPair with empty secret key: expected error %q, got %v", expectedErr, err)
	}
}

// createExpiredToken generates a token that expires very quickly.
func createExpiredToken(userID int64, secretKey string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)), // Expired 1 minute ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-5 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func TestValidateToken(t *testing.T) {
	// Valid token
	validToken, err := GenerateAccessToken(testUserID, testSecret)
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err)
	}

	// Expired token
	expiredToken, err := createExpiredToken(testUserID, testSecret)
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	// Token with wrong secret
	wrongSecret := "wrongsecret"
	tokenWithWrongSecret, err := GenerateAccessToken(testUserID, wrongSecret)
	if err != nil {
		t.Fatalf("Failed to generate token with wrong secret: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		secretKey   string
		expectError bool
		expectedMsg string
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			secretKey:   testSecret,
			expectError: false,
		},
		{
			name:        "Expired Token",
			tokenString: expiredToken,
			secretKey:   testSecret,
			expectError: true,
			expectedMsg: "token expired",
		},
		{
			name:        "Invalid Signature",
			tokenString: tokenWithWrongSecret,
			secretKey:   testSecret,
			expectError: true,
			expectedMsg: "invalid token signature",
		},
		{
			name:        "Malformed Token",
			tokenString: "malformed.token.string",
			secretKey:   testSecret,
			expectError: true,
			expectedMsg: "token is malformed: could not base64 decode header: illegal base64 data at input byte 8",
		},
		{
			name:        "Empty Secret Key",
			tokenString: validToken,
			secretKey:   "",
			expectError: true,
			expectedMsg: "secret key cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.tokenString, tt.secretKey)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateToken() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if tt.expectError {
				if err == nil || err.Error() != tt.expectedMsg {
					t.Errorf("ValidateToken() got error message = %q, want %q", err.Error(), tt.expectedMsg)
				}
			} else {
				if claims.UserID != testUserID {
					t.Errorf("ValidateToken() got UserID = %d, want %d", claims.UserID, testUserID)
				}
			}
		})
	}
}

func TestExtractUserID(t *testing.T) {
	// Valid token
	validToken, err := GenerateAccessToken(testUserID, testSecret)
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err)
	}

	// Expired token
	expiredToken, err := createExpiredToken(testUserID, testSecret)
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		secretKey   string
		expectError bool
		expectedID  int64
		expectedMsg string
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			secretKey:   testSecret,
			expectError: false,
			expectedID:  testUserID,
		},
		{
			name:        "Expired Token",
			tokenString: expiredToken,
			secretKey:   testSecret,
			expectError: true,
			expectedID:  0,
			expectedMsg: "token expired",
		},
		{
			name:        "Invalid Token (wrong secret)",
			tokenString: validToken,
			secretKey:   "anothersecret",
			expectError: true,
			expectedID:  0,
			expectedMsg: "invalid token signature",
		},
		{
			name:        "Empty Secret Key",
			tokenString: validToken,
			secretKey:   "",
			expectError: true,
			expectedID:  0,
			expectedMsg: "secret key cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := ExtractUserID(tt.tokenString, tt.secretKey)
			if (err != nil) != tt.expectError {
				t.Errorf("ExtractUserID() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if tt.expectError {
				if err == nil || err.Error() != tt.expectedMsg {
					t.Errorf("ExtractUserID() got error message = %q, want %q", err.Error(), tt.expectedMsg)
				}
			} else {
				if userID != tt.expectedID {
					t.Errorf("ExtractUserID() got UserID = %d, want %d", userID, tt.expectedID)
				}
			}
		})
	}
}
