package my_struct

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID uint32
	jwt.StandardClaims
}
