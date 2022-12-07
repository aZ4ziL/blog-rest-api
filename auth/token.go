package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Declare the secret key
var SecretKey = []byte("mysecretkey")

// Credential
// will declare the credentials for the user.
type Credential struct {
	ID          uint      `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	IsSuperuser bool      `json:"is_superuser"`
	IsStaff     bool      `json:"is_staff"`
	IsActive    bool      `json:"is_active"`
	LastLogin   time.Time `json:"last_login"`
	DateJoined  time.Time `json:"date_joined"`
}

// Claims
type Claims struct {
	Credential
	jwt.RegisteredClaims
}

// GenerateNewToken
// generate new token
func GenerateNewToken(cred Credential) (token string, err error) {
	// set expiration time for tokens.
	expirationTime := time.Now().Add(24 * time.Hour) // 24 Hour

	claims := &Claims{
		Credential: cred,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// get token string
	token, err = tokenJWT.SignedString(SecretKey)
	return
}

// ReadAndVerifyToken
// decode the token and convert to claims
func ReadAndVerifyToken(token string) (Claims, error) {
	var claims Claims

	tokenJWT, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return Claims{}, err
	}

	// Check token is valid
	if !tokenJWT.Valid {
		return Claims{}, errors.New("error: Your token is not valid")
	}
	return claims, err
}
