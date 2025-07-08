package middlewares

import (
	"github.com/alpha-omega-corp/cloud/api/clients/user"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type AuthMiddleware struct {
	userClient user.Client
}

func NewAuthMiddleware(client user.Client) *AuthMiddleware {
	return &AuthMiddleware{
		userClient: client,
	}
}

func (m *AuthMiddleware) Auth(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		err := m.userClient.Validate(w, req)
		if err != nil {
			return err
		}

		return next(w, req)
	}
}
