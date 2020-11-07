package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/divisi-developer-poros/poros-web-backend/config"
	"github.com/twinj/uuid"
)

// Variable Used
var (
	EnvironmentToken config.TokenENV
	dbModel          config.DBModel
	connection       = dbModel.PostgreConn()
)

// GenerateToken ... Generate JWT Token
func (jt *JWTToken) GenerateToken(userName string, userType int) (string, error) {
	tokenID := uuid.NewV1().String()
	expiresAt := time.Now().Add(time.Hour * 48).Unix()

	result, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTToken{
		Username: userName,
		Usertype: userType,
		ID:       tokenID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "poros",
			IssuedAt:  time.Now().Unix(),
		},
	}).SignedString([]byte(EnvironmentToken.JWTSecret))
	if err != nil {
		return "", err
	}
	accessToken := AccessToken{
		ID:        tokenID,
		ExpiresAt: expiresAt,
	}
	if err := connection.Create(&accessToken).Error; err != nil {
		return "", err
	}
	return result, nil
}

// TokenValidation ... Validate JWT Token
func (jt *JWTToken) TokenValidation(encodedToken string) (*jwt.Token, error) {
	result, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, errors.New("invalid token")
		}
		return []byte(EnvironmentToken.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := result.Claims.(jwt.MapClaims)
	if ok {
		tokenID, ok := claims["id"].(string)
		if !ok {
			return nil, err
		}
		var accessToken AccessToken
		if err := connection.Where("id = ?", tokenID).Find(&accessToken).Error; err != nil {
			return nil, err
		}
		if accessToken.ID == "" {
			return nil, errors.New("logged out")
		}
	}
	return result, err
}

// DeleteToken ... Delete JWT Token from DB
func (jt *JWTToken) DeleteToken(id string) error {
	var accessToken AccessToken
	if err := connection.Where("id = ?", id).Find(&accessToken).Error; err != nil {
		return err
	}
	if err := connection.Delete(&accessToken).Error; err != nil {
		return err
	}
	return nil
}
