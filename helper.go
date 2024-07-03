package crmhazar_pkg_http

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

func TokenClaims(token, secretKey string) (jwt.MapClaims, error) {
	decoded, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)

	if !ok {
		// TODO tokenin omrini test etmeli
		// TODO son seretmeli
		return nil, errors.New("internal server error")
	}

	return claims, nil
}
