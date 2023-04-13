package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID   uint32
	Name string
	jwt.StandardClaims
}
