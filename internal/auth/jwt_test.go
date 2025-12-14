package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	testSecretKey = "supersecretkey"
	testAccountID = int64(123)
	testLenderID  = int64(456)
)

// parseToken parses and validates a JWT token string.
func parseToken(t *testing.T, tokenString string, secretKey string) *Claims {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Fatalf("Token is not valid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatalf("Failed to get claims from token")
	}

	return claims
}

func TestGenerateAccessToken(t *testing.T) {
	tokenString, err := GenerateAccessToken(testAccountID, testLenderID, testSecretKey)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}
	if tokenString == "" {
		t.Fatal("Generated access token is empty")
	}

	claims := parseToken(t, tokenString, testSecretKey)

	if claims.AccountID != testAccountID {
		t.Errorf("Expected AccountID %d, got %d", testAccountID, claims.AccountID)
	}
	if claims.LenderID != testLenderID {
		t.Errorf("Expected LenderID %d, got %d", testLenderID, claims.LenderID)
	}

	// Check expiration time (allow for a small margin of error)
	expectedExpiry := time.Now().Add(AccessTokenDuration)
	if claims.ExpiresAt == nil {
		t.Fatal("AccessToken claims.ExpiresAt is nil")
	}
	if claims.ExpiresAt.Time.Before(expectedExpiry.Add(-1*time.Minute)) || claims.ExpiresAt.Time.After(expectedExpiry.Add(1*time.Minute)) {
		t.Errorf("AccessToken expiration time is not within expected range. Expected around %v, got %v", expectedExpiry, claims.ExpiresAt.Time)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tokenString, err := GenerateRefreshToken(testAccountID, testLenderID, testSecretKey)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if tokenString == "" {
		t.Fatal("Generated refresh token is empty")
	}

	claims := parseToken(t, tokenString, testSecretKey)

	if claims.AccountID != testAccountID {
		t.Errorf("Expected AccountID %d, got %d", testAccountID, claims.AccountID)
	}
	if claims.LenderID != testLenderID {
		t.Errorf("Expected LenderID %d, got %d", testLenderID, claims.LenderID)
	}

	// Check expiration time (allow for a small margin of error)
	expectedExpiry := time.Now().Add(RefreshTokenDuration)
	if claims.ExpiresAt == nil {
		t.Fatal("RefreshToken claims.ExpiresAt is nil")
	}
	if claims.ExpiresAt.Time.Before(expectedExpiry.Add(-1*time.Minute)) || claims.ExpiresAt.Time.After(expectedExpiry.Add(1*time.Minute)) {
		t.Errorf("RefreshToken expiration time is not within expected range. Expected around %v, got %v", expectedExpiry, claims.ExpiresAt.Time)
	}
}

func TestGenerateTokenPair(t *testing.T) {
	tokenPair, err := GenerateTokenPair(testAccountID, testLenderID, testSecretKey)
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}
	if tokenPair == nil {
		t.Fatal("Generated token pair is nil")
	}
	if tokenPair.AccessToken == "" {
		t.Fatal("Generated access token in pair is empty")
	}
	if tokenPair.RefreshToken == "" {
		t.Fatal("Generated refresh token in pair is empty")
	}

	// Validate Access Token
	accessClaims := parseToken(t, tokenPair.AccessToken, testSecretKey)
	if accessClaims.AccountID != testAccountID || accessClaims.LenderID != testLenderID {
		t.Errorf("Access Token claims mismatch: AccountID %d/%d, LenderID %d/%d",
			accessClaims.AccountID, testAccountID, accessClaims.LenderID, testLenderID)
	}
	expectedAccessExpiry := time.Now().Add(AccessTokenDuration)
	if accessClaims.ExpiresAt == nil {
		t.Fatal("AccessToken claims.ExpiresAt is nil")
	}
	if accessClaims.ExpiresAt.Time.Before(expectedAccessExpiry.Add(-1*time.Minute)) || accessClaims.ExpiresAt.Time.After(expectedAccessExpiry.Add(1*time.Minute)) {
		t.Errorf("Access Token expiration time is not within expected range. Expected around %v, got %v", expectedAccessExpiry, accessClaims.ExpiresAt.Time)
	}

	// Validate Refresh Token
	refreshClaims := parseToken(t, tokenPair.RefreshToken, testSecretKey)
	if refreshClaims.AccountID != testAccountID || refreshClaims.LenderID != testLenderID {
		t.Errorf("Refresh Token claims mismatch: AccountID %d/%d, LenderID %d/%d",
			refreshClaims.AccountID, testAccountID, refreshClaims.LenderID, testLenderID)
	}
	expectedRefreshExpiry := time.Now().Add(RefreshTokenDuration)
	if refreshClaims.ExpiresAt == nil {
		t.Fatal("RefreshToken claims.ExpiresAt is nil")
	}
	if refreshClaims.ExpiresAt.Time.Before(expectedRefreshExpiry.Add(-1*time.Minute)) || refreshClaims.ExpiresAt.Time.After(expectedRefreshExpiry.Add(1*time.Minute)) {
		t.Errorf("Refresh Token expiration time is not within expected range. Expected around %v, got %v", expectedRefreshExpiry, refreshClaims.ExpiresAt.Time)
	}
}

func TestValidateToken_Valid(t *testing.T) {
	tokenString, err := GenerateAccessToken(testAccountID, testLenderID, testSecretKey)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateToken(tokenString, testSecretKey)
	if err != nil {
		t.Fatalf("ValidateToken failed for valid token: %v", err)
	}

	if claims.AccountID != testAccountID {
		t.Errorf("Expected AccountID %d, got %d", testAccountID, claims.AccountID)
	}
	if claims.LenderID != testLenderID {
		t.Errorf("Expected LenderID %d, got %d", testLenderID, claims.LenderID)
	}
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	tokenString, err := GenerateAccessToken(testAccountID, testLenderID, testSecretKey)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = ValidateToken(tokenString, "wrongsecretkey")
	if err == nil {
		t.Fatal("ValidateToken unexpectedly succeeded with wrong secret key")
	}
	if !errors.Is(err, jwt.ErrSignatureInvalid) {
		t.Errorf("Expected signature invalid error, got: %v", err)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	// Create a token with a very short expiry time
	expiredClaims := Claims{
		AccountID: testAccountID,
		LenderID:  testLenderID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(100 * time.Millisecond)), // Short expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, err := token.SignedString([]byte(testSecretKey))
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}

	time.Sleep(150 * time.Millisecond) // Wait for the token to expire

	_, err = ValidateToken(tokenString, testSecretKey)
	if err == nil {
		t.Fatal("ValidateToken unexpectedly succeeded for expired token")
	}
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("Expected ErrTokenExpired, got: %v", err)
	}
}

func TestExtractAccountID(t *testing.T) {
	tokenString, err := GenerateAccessToken(testAccountID, testLenderID, testSecretKey)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	extractedID, err := ExtractAccountID(tokenString, testSecretKey)
	if err != nil {
		t.Fatalf("ExtractAccountID failed: %v", err)
	}
	if extractedID != testAccountID {
		t.Errorf("Expected extracted AccountID %d, got %d", testAccountID, extractedID)
	}

	// Test with invalid token
	_, err = ExtractAccountID("invalid.token.string", testSecretKey)
	if err == nil {
		t.Fatal("ExtractAccountID unexpectedly succeeded with invalid token string")
	}
}

func TestExtractLenderID(t *testing.T) {
	tokenString, err := GenerateAccessToken(testAccountID, testLenderID, testSecretKey)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	extractedID, err := ExtractLenderID(tokenString, testSecretKey)
	if err != nil {
		t.Fatalf("ExtractLenderID failed: %v", err)
	}
	if extractedID != testLenderID {
		t.Errorf("Expected extracted LenderID %d, got %d", testLenderID, extractedID)
	}

	// Test with invalid token
	_, err = ExtractLenderID("invalid.token.string", testSecretKey)
	if err == nil {
		t.Fatal("ExtractLenderID unexpectedly succeeded with invalid token string")
	}
}
