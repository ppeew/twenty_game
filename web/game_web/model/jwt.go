package model

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID int
	jwt.StandardClaims
}
