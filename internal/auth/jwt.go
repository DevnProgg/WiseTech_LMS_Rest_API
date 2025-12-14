package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	AccountID int64
	LenderID  int64
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new access token for the given account and lender IDs.
func GenerateAccessToken(accountID, lenderID int64, secretKey string) (string, error) {
	claims := Claims{
		AccountID: accountID,
		LenderID:  lenderID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// GenerateRefreshToken creates a new refresh token for the given account and lender IDs.
func GenerateRefreshToken(accountID, lenderID int64, secretKey string) (string, error) {
	claims := Claims{
		AccountID: accountID,
		LenderID:  lenderID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// GenerateTokenPair generates both an access token and a refresh token.
func GenerateTokenPair(accountID, lenderID int64, secretKey string) (*TokenPair, error) {
	accessToken, err := GenerateAccessToken(accountID, lenderID, secretKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(accountID, lenderID, secretKey)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken parses and validates a JWT token string, returning its claims if valid.
func ValidateToken(tokenString, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// ExtractAccountID extracts the AccountID from a validated token.
func ExtractAccountID(tokenString, secretKey string) (int64, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return 0, err
	}
	return claims.AccountID, nil
}

// ExtractLenderID extracts the LenderID from a validated token.
func ExtractLenderID(tokenString, secretKey string) (int64, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return 0, err
	}
	return claims.LenderID, nil
}