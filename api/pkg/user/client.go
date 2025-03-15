package user

import (
	"fmt"
	"github.com/alpha-omega-corp/services/types"
	"github.com/alpha-omega-corp/user-svc/proto"
	"github.com/uptrace/bunrouter"
	"google.golang.org/grpc"
	"net/http"
)

type Client interface {
	Self() proto.UserServiceClient
	Login(w http.ResponseWriter, req bunrouter.Request) error
	Register(w http.ResponseWriter, req bunrouter.Request) error
	GetUsers(w http.ResponseWriter, req bunrouter.Request) error
	CreateUser(w http.ResponseWriter, req bunrouter.Request) error
	UpdateUser(w http.ResponseWriter, req bunrouter.Request) error
	DeleteUser(w http.ResponseWriter, req bunrouter.Request) error
	AssignUser(w http.ResponseWriter, req bunrouter.Request) error
	GetServices(w http.ResponseWriter, req bunrouter.Request) error
	GetServicePermissions(w http.ResponseWriter, req bunrouter.Request) error
	CreateServicePermissions(w http.ResponseWriter, req bunrouter.Request) error
	GetUserPermissions(w http.ResponseWriter, req bunrouter.Request) error
	GetRoles(w http.ResponseWriter, req bunrouter.Request) error
	CreateRole(w http.ResponseWriter, req bunrouter.Request) error
	GetTest(w http.ResponseWriter, req bunrouter.Request) error
}

type userClient struct {
	Client
	client proto.UserServiceClient
}

func NewClient(c types.ConfigHost) Client {
	conn, err := grpc.Dial(c.Url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return &userClient{client: proto.NewUserServiceClient(conn)}
}

func (svc *userClient) Self() proto.UserServiceClient {
	return svc.client
}
func (svc *userClient) Login(w http.ResponseWriter, req bunrouter.Request) error {
	return LoginHandler(w, req, svc.client)
}
func (svc *userClient) Register(w http.ResponseWriter, req bunrouter.Request) error {
	return RegisterHandler(w, req, svc.client)
}
func (svc *userClient) GetUsers(w http.ResponseWriter, req bunrouter.Request) error {
	return GetUsersHandler(w, req, svc.client)
}
func (svc *userClient) CreateUser(w http.ResponseWriter, req bunrouter.Request) error {
	return CreateUserHandler(w, req, svc.client)
}
func (svc *userClient) UpdateUser(w http.ResponseWriter, req bunrouter.Request) error {
	return UpdateUserHandler(w, req, svc.client)
}
func (svc *userClient) DeleteUser(w http.ResponseWriter, req bunrouter.Request) error {
	return DeleteUserHandler(w, req, svc.client)
}
func (svc *userClient) AssignUser(w http.ResponseWriter, req bunrouter.Request) error {
	return AssignUserHandler(w, req, svc.client)
}
func (svc *userClient) GetUserPermissions(w http.ResponseWriter, req bunrouter.Request) error {
	return GetUserPermissionsHandler(w, req, svc.client)
}
func (svc *userClient) GetServicePermissions(w http.ResponseWriter, req bunrouter.Request) error {
	return GetServicePermissionsHandler(w, req, svc.client)
}
func (svc *userClient) CreatePermission(w http.ResponseWriter, req bunrouter.Request) error {
	return CreatePermissionHandler(w, req, svc.client)
}
func (svc *userClient) GetServices(w http.ResponseWriter, req bunrouter.Request) error {
	return GetServices(w, req, svc.client)
}
func (svc *userClient) GetRoles(w http.ResponseWriter, req bunrouter.Request) error {
	return GetRolesHandler(w, req, svc.client)
}
func (svc *userClient) CreateRole(w http.ResponseWriter, req bunrouter.Request) error {
	return CreateRoleHandler(w, req, svc.client)
}
func (svc *userClient) GetTest(w http.ResponseWriter, req bunrouter.Request) error {
	return GetTestHandler(w, req, svc.client)
}
