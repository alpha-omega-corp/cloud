package user

import (
	_ "github.com/spf13/viper/remote"
	"github.com/uptrace/bunrouter"
)

func RegisterClient(client Client, r *bunrouter.Router) Client {

	r.POST("/auth/login", client.Login)
	r.POST("/auth/register", client.Register)

	r.GET("/auth/roles", client.GetRoles)
	r.POST("/auth/roles", client.CreateRole)
	r.GET("/auth/services", client.GetServices)
	r.GET("/auth/services/:id/permissions", client.GetServicePermissions)
	r.POST("/auth/services/permissions", client.CreateServicePermissions)

	r.GET("/users", client.GetUsers)
	r.POST("/users", client.CreateUser)
	r.PUT("/users/:id", client.UpdateUser)
	r.DELETE("/users/:id", client.DeleteUser)
	r.POST("/users/roles", client.AssignUser)
	r.GET("/users/:id/permissions", client.GetUserPermissions)

	return client
}
