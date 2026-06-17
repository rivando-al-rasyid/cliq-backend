package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AccessTokenExpiry is the lifetime of an access JWT.
const AccessTokenExpiry = 15 * time.Minute

// ResetTokenExpiry is the lifetime of a short-lived password-reset JWT.
const ResetTokenExpiry = 10 * time.Minute

// // AccessTokenSubject is the JWT "sub" claim used for normal access token.
// const AccessTokenSubject = "access"

// ResetTokenSubject is the JWT "sub" claim used exclusively for password-reset JWTs.
const ResetTokenSubject = "password-reset"

type Claims struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

func NewClaims(id uuid.UUID, email string) *Claims {
	return &Claims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			Subject:   AccessTokenSubject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func NewResetClaims(id uuid.UUID, email string) *Claims {
	return &Claims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			Subject:   ResetTokenSubject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ResetTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func (c *Claims) GenJWT() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("missing jwt secret")
	}

	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		return "", errors.New("missing jwt issuer")
	}

	if c.Issuer == "" {
		c.Issuer = jwtIssuer
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (c *Claims) VerifyJWT(rawToken string) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("missing jwt secret")
	}

	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		return errors.New("missing jwt issuer")
	}

	token, err := jwt.ParseWithClaims(rawToken, c, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrTokenInvalidClaims
	}

	if c.Issuer != jwtIssuer {
		return jwt.ErrTokenInvalidIssuer
	}

	return nil
}
