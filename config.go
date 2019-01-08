package main

import (
	"errors"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Config of the api application
type Config struct {
	// HTTP port to listen
	HTTPPort int

	TokenConfig *TokenConfig
}

// TokenConfig configuration for the Token service
type TokenConfig struct {
	// Key is the key to encrypt/decrypt the token
	Key string

	// ExpirationTime duration in seconds of the token validity
	ExpirationTime int

	// Issuer identify the principal that issued the token
	Issuer string

	// Audience identify the recipients that the token is intended for
	Audience string
}

// GetConfig return the Config loaded from environment variables
func GetConfig() *Config {
	return &Config{
		HTTPPort:    getPort(),
		TokenConfig: getTokenConfig(),
	}
}

func getPort() (pInt int) {
	var (
		p   string
		err error
	)

	if *portFlag != "" {
		p = *portFlag
	} else if p = os.Getenv("JWT_PORT"); p == "" {
		err := errors.New("Missing HTTP Port")
		log.Fatal(err)
	}

	pInt, err = strconv.Atoi(p)
	if err != nil {
		err = errors.New("Invalid port number")
		log.Fatal(err)
	}

	return
}

func getTokenConfig() (tc *TokenConfig) {
	var err error

	tc = &TokenConfig{}

	if *tokenKeyFlag != "" {
		tc.Key = *tokenKeyFlag
	} else if tc.Key = os.Getenv("JWT_KEY"); tc.Key == "" {
		err = errors.New("Missing key")
		log.Fatal(err)
	}

	if *tokenIssuerFlag != "" {
		tc.Issuer = *tokenIssuerFlag
	} else {
		tc.Issuer = os.Getenv("JWT_ISSUER")
	}

	if *tokenAudienceFlag != "" {
		tc.Audience = *tokenAudienceFlag
	} else {
		tc.Audience = os.Getenv("JWT_AUDIENCE")
	}

	var d string
	if *tokenDurationFlag != "" {
		d = *tokenDurationFlag
	} else if d = os.Getenv("JWT_DURATION"); d == "" {
		err = errors.New("Missing token duration")
		log.Fatal(err)
	}

	tc.ExpirationTime, err = strconv.Atoi(d)
	if err != nil {
		err = errors.New("Invalid token duration")
		log.Fatal(err)
	}

	return
}
