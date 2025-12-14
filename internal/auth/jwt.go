package auth

import (

	"errors"

	"time"



	"github.com/golang-jwt/jwt/v5"

)

// AccessTokenDuration defines how long an access token is valid.

const AccessTokenDuration = 15 * time.Minute



// RefreshTokenDuration defines how long a refresh token is valid.

const RefreshTokenDuration = 7 * 24 * time.Hour



// TokenPair holds both the access and refresh tokens.

type TokenPair struct {

	AccessToken  string

	RefreshToken string

}



// Claims represents the JWT claims, embedding jwt.RegisteredClaims and adding UserID.

type Claims struct {

	UserID int64 `json:"user_id"`

	jwt.RegisteredClaims

}



// GenerateAccessToken generates a new access token for a given user ID.

func GenerateAccessToken(userID int64, secretKey string) (string, error) {

	if secretKey == "" {

		return "", errors.New("secret key cannot be empty")

	}
	claims := Claims{
		UserID: userID,
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

// GenerateRefreshToken generates a new refresh token for a given user ID.
func GenerateRefreshToken(userID int64, secretKey string) (string, error) {
	if secretKey == "" {
		return "", errors.New("secret key cannot be empty")
	}
	claims := Claims{
		UserID: userID,
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
func GenerateTokenPair(userID int64, secretKey string) (*TokenPair, error) {
	accessToken, err := GenerateAccessToken(userID, secretKey)
	if err != nil {
		return nil, err
	}
	refreshToken, err := GenerateRefreshToken(userID, secretKey)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken parses and validates a JWT token string.
func ValidateToken(tokenString, secretKey string) (*Claims, error) {
	if secretKey == "" {
		return nil, errors.New("secret key cannot be empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	}, jwt.WithLeeway(5*time.Second)) // Add a small leeway for clock skew

	if err != nil {
		// Handle specific JWT errors
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err // Other parsing errors
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractUserID validates the token and extracts the UserID from its claims.
func ExtractUserID(tokenString, secretKey string) (int64, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
