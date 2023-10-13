package model

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID uint32
	jwt.StandardClaims
}
