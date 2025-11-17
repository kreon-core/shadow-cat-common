package middlware

import (
	"net/http"

	"github.com/go-chi/cors"

	tul "github.com/kreon-core/shadow-cat-common"
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
		AllowedOrigins: tul.OrElse(cfg.AllowedOrigins, []string{"https://*", "http://*"}),
		AllowedMethods: tul.OrElse(cfg.AllowedMethods, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		AllowedHeaders: tul.OrElse(cfg.AllowedHeaders, []string{
			"Origin",
			"Accept",
			"Content-Type",
			"Authorization",
			"X-Real-IP",
			"X-Request-ID",
		}),
		ExposedHeaders:   tul.OrElse(cfg.ExposedHeaders, []string{"Content-Length"}),
		AllowCredentials: tul.OrElse(cfg.AllowCredentials, false),
		MaxAge:           tul.OrElse(cfg.MaxAge, maxAge),
	})
}
