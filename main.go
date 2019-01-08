package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Route describe a route
type Route struct {
	Name    string
	Path    string
	Method  string
	Handler http.Handler
	Public  bool
}

var (
	portFlag          = flag.String("port", "", "The port to listen on")
	tokenKeyFlag      = flag.String("key", "", "The token key")
	tokenIssuerFlag   = flag.String("issuer", "", "The token issuer")
	tokenDurationFlag = flag.String("duration", "", "The token validity duration (in seconds)")
	tokenAudienceFlag = flag.String("audience", "", "The token audience")
)

func main() {

	conf := GetConfig()

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000",
		FullTimestamp:   true,
	})

	jwtServ := &JWTService{
		TokenKey:            conf.TokenConfig.Key,
		TokenIssuer:         conf.TokenConfig.Issuer,
		TokenExpirationTime: time.Duration(conf.TokenConfig.ExpirationTime) * time.Second,
		TokenAudience:       conf.TokenConfig.Audience,
	}

	routes := []Route{
		{
			Name:   `tokenNew`,
			Method: http.MethodPost,
			Handler: TokenNewHandler{
				JWTService: jwtServ,
			},
			Path: `/tokens`,
		},
		{
			Name:   `tokenCheck`,
			Method: http.MethodGet,
			Handler: TokenCheckHandler{
				JWTService: jwtServ,
			},
			Path: `/tokens/check`,
		},
		{
			Name:   `tokenDecode`,
			Method: http.MethodGet,
			Handler: TokenDecodeHandler{
				JWTService: jwtServ,
			},
			Path: `/tokens/decode`,
		},
	}

	router := mux.NewRouter()

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.Handler)
	}

	log.Info(fmt.Sprintf("Listening on %d", conf.HTTPPort))
	log.SetReportCaller(true)

	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.HTTPPort), router)
	if err != nil {
		log.Fatal(err)
	}
}
