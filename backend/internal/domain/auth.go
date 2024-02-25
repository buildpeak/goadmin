package domain

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
