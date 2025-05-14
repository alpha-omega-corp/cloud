package middlewares

import (
	"github.com/alpha-omega-corp/cloud/api/clients/user"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/proto"
	"github.com/uptrace/bunrouter"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	service user.Client
}

func NewAuthMiddleware(userService user.Client) *AuthMiddleware {
	return &AuthMiddleware{
		service: userService,
	}
}

func (middleware *AuthMiddleware) Auth(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		authHeader := req.Header.Get("Authorization")
		token := strings.Split(authHeader, "Bearer ")[1]

		_, err := middleware.service.Self().Validate(req.Context(), &proto.ValidateRequest{
			Token: token,
		})

		if err != nil {
			return err
		}

		return next(w, req)
	}
}
