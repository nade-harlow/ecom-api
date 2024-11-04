package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/nade-harlow/ecom-api/internal/config"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	secretKey = config.AppConfig.JwtSecret
)

type DecodedUser struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
}

func GenerateJwt(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, error := token.SignedString([]byte(secretKey))
	if error != nil {
		fmt.Println(error)
	}
	return tokenString
}

func GetJWTClaims(userId uuid.UUID, role string) jwt.MapClaims {
	now := time.Now()
	issuedTime := now.Unix()
	accessExpires := now.Add(time.Minute * 30).Unix()

	return jwt.MapClaims{
		"userId": userId,
		"iat":    issuedTime,
		"exp":    accessExpires,
		"role":   role,
	}
}

func ValidateJwtAuthenticity(token string) (jwt.Claims, error) {
	validToken, tokenError := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error validating Token")
		}
		return []byte(secretKey), nil
	})

	if tokenError != nil && !validToken.Valid {
		return nil, fmt.Errorf("invalid authorization token")
	}

	return validToken.Claims, nil
}

func GetLoginUser(context *gin.Context) (DecodedUser, error) {
	var user DecodedUser
	decodeToken, exist := context.Get("decodedToken")

	if !exist {
		return DecodedUser{}, errors.New("invalid access to route")
	}

	err := json.Unmarshal(decodeToken.([]byte), &user)
	if err != nil {
		fmt.Println("error decoding token", err)
		return DecodedUser{}, err
	}

	return user, nil
}
