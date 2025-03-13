package user

import (
	_ "github.com/spf13/viper/remote"
	"github.com/uptrace/bunrouter"
)

func RegisterClient(svc Client, r *bunrouter.Router) Client {

	r.POST("/login", svc.Login)
	r.POST("/register", svc.Register)
	r.GET("/roles", svc.GetRoles)
	r.POST("/role", svc.CreateRole)
	r.GET("/users", svc.GetUsers)
	r.POST("/user", svc.CreateUser)
	r.PUT("/user/:id", svc.UpdateUser)
	r.DELETE("/user/:id", svc.DeleteUser)
	r.POST("/user/assign", svc.AssignUser)
	r.GET("/user/:id/permissions", svc.GetUserPermissions)
	r.GET("/services", svc.GetServices)
	r.GET("/service/:serviceId/permissions", svc.GetServicePermissions)
	r.POST("/service/permissions", svc.CreateServicePermissions)
	r.GET("/user/test", svc.GetTest)

	return svc
}
