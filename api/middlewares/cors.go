package middlewares

import (
	"github.com/rs/cors"
	"github.com/uptrace/bunrouter"
)

func NewCorsMiddleware() bunrouter.MiddlewareFunc {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4000"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	})

	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return bunrouter.HTTPHandler(corsHandler.Handler(next))
	}
}
