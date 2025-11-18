package mwc

import (
	"net/http"

	"github.com/go-chi/cors"

	"github.com/kreon-core/shadow-cat-common/utlc"
)

const maxAge = 300

type CORSConfig struct {
	AllowedOrigins   *[]string
	AllowOriginFunc  *func(r *http.Request, origin string) bool
	AllowedMethods   *[]string
	AllowedHeaders   *[]string
	ExposedHeaders   *[]string
	AllowCredentials *bool
	MaxAge           *int
}

func CORS(cfg *CORSConfig) func(next http.Handler) http.Handler {
	if cfg == nil {
		cfg = &CORSConfig{}
	}

	corsOptions := &cors.Options{
		AllowedMethods: utlc.OrElse(cfg.AllowedMethods, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		AllowedHeaders: utlc.OrElse(cfg.AllowedHeaders, []string{
			"Origin",
			"Accept",
			"Content-Type",
			"Authorization",
			"X-Real-IP",
			"X-Request-ID",
		}),
		ExposedHeaders:   utlc.OrElse(cfg.ExposedHeaders, []string{"Content-Length"}),
		AllowCredentials: utlc.OrElse(cfg.AllowCredentials, false),
		MaxAge:           utlc.OrElse(cfg.MaxAge, maxAge),
	}

	if cfg.AllowOriginFunc != nil {
		corsOptions.AllowOriginFunc = *cfg.AllowOriginFunc
	} else {
		corsOptions.AllowedOrigins = utlc.OrElse(cfg.AllowedOrigins, []string{"https://*", "http://*"})
	}

	return cors.Handler(*corsOptions)
}
