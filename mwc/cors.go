package mwc

import (
	"net/http"

	"github.com/go-chi/cors"

	"github.com/kreon-core/shadow-cat-common/ultc"
)

const maxAge = 300

type CORSConfig struct {
	AllowedOrigins   *[]string
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
	return cors.Handler(cors.Options{
		AllowedOrigins: ultc.OrElse(cfg.AllowedOrigins, []string{"https://*", "http://*"}),
		AllowedMethods: ultc.OrElse(cfg.AllowedMethods, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		AllowedHeaders: ultc.OrElse(cfg.AllowedHeaders, []string{
			"Origin",
			"Accept",
			"Content-Type",
			"Authorization",
			"X-Real-IP",
			"X-Request-ID",
		}),
		ExposedHeaders:   ultc.OrElse(cfg.ExposedHeaders, []string{"Content-Length"}),
		AllowCredentials: ultc.OrElse(cfg.AllowCredentials, false),
		MaxAge:           ultc.OrElse(cfg.MaxAge, maxAge),
	})
}
