package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"strconv"
)


// need to be loaded from env
var jwtSecret string = "xxxxxxxxxxxxxxxxxxxx"

func createClaimsForUser(id uint64) *jwt.RegisteredClaims{

	return &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   "Authentication",
		ID:        strconv.FormatUint((id), 10),
	}
}


func CreateToken(id uint64) (string, error) {
	claims := createClaimsForUser(id)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) 
	res , err := token.SignedString([]byte(jwtSecret))

	if err != nil{
		return "",  err
	}

	return res, nil
}


func VerifyToken(tokenString string) (string , bool)  {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil{
		return "", false
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if ok {
		return claims.ID, true
	}
	return "", false
}
