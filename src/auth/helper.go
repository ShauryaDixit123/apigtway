package auth

import (
	"apigtway/src/dtos"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request, typetkn string) (*jwt.Token, error) {
	var typeofKey string
	switch typetkn {
	case "access":
		typeofKey = "ACCESS_SECRET_KEY"
	case "refresh":
		typeofKey = "REFRESH_SECRET_KEY"
	}
	tokenString := extractToken(r)
	token, err := ParseToken(tokenString, typeofKey)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ParseToken(tkn string, key string) (*jwt.Token, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(key)), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request, typetkn string) error {
	token, err := VerifyToken(r, typetkn)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractMetaData(r *http.Request, typetkn string) (*dtos.AccessDetails, error) {
	var struuid string
	token, err := VerifyToken(r, typetkn)
	if err != nil {
		return nil, err
	}
	switch typetkn {
	case "access":
		struuid = "access_uuid"
	case "refresh":
		struuid = "refresh_uuid"
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims[struuid].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &dtos.AccessDetails{
			AccessUUID: accessUuid,
			UserID:     userId,
		}, nil
	}
	return nil, err
}
