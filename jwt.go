package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// JWTService handling token signature and validation
type JWTService struct {
	TokenKey            string
	TokenIssuer         string
	TokenAudience       string
	TokenExpirationTime time.Duration
}

// NewTokenString generates a new Token
func (s *JWTService) NewTokenString(data map[string]interface{}) (tokenString string, err error) {

	if s.TokenIssuer != `` {
		data["iss"] = s.TokenIssuer
	}

	if s.TokenAudience != `` {
		data["aud"] = s.TokenAudience
	}

	data["exp"] = time.Now().Add(s.TokenExpirationTime).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(data))

	tokenString, err = token.SignedString([]byte(s.TokenKey))
	if err != nil {
		log.Error(err)
		return tokenString, err
	}

	return
}

// ValidateToken validates the token signature and returns the embedded data
func (s *JWTService) ValidateToken(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.TokenKey), nil
	})

	if token != nil && token.Valid {
		ts := strings.Split(tokenString, `.`)
		if len(ts) != 3 {
			err = errors.New("token contains an invalid number of segments")
			log.Error(err)
			return nil, err
		}

		d, err := jwt.DecodeSegment(ts[1])
		if err != nil {
			log.Error(err)
			return nil, err
		}

		var data map[string]interface{}
		if err = json.Unmarshal(d, &data); err != nil {
			log.Error(err)
			return nil, err
		}

		return data, nil
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, fmt.Errorf("String provided is not a tokenString : %s", tokenString)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// fmt.Println("Expired")
			return nil, nil
		} else {
			log.Error(err)
			return nil, err
		}
	}

	log.Error(err)
	return nil, err

}
