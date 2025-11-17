package tul

import (
	"net/http"

	"github.com/go-chi/cors"
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
	return cors.Handler(cors.Options{
		AllowedOrigins: OrElse(cfg.AllowedOrigins, []string{"https://*", "http://*"}),
		AllowedMethods: OrElse(cfg.AllowedMethods, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		AllowedHeaders: OrElse(cfg.AllowedHeaders, []string{
			"Origin",
			"Accept",
			"Content-Type",
			"Authorization",
			"X-Real-IP",
			"X-Request-ID",
		}),
		ExposedHeaders:   OrElse(cfg.ExposedHeaders, []string{"Content-Length"}),
		AllowCredentials: OrElse(cfg.AllowCredentials, false),
		MaxAge:           OrElse(cfg.MaxAge, maxAge),
	})
}
