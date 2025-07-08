package user

import (
	"fmt"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/proto"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/alpha-omega-corp/cloud/core/httputils"
	"github.com/uptrace/bunrouter"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strings"
)

type Client interface {
	Login(w http.ResponseWriter, req bunrouter.Request) error
	Validate(w http.ResponseWriter, req bunrouter.Request) error
	Register(w http.ResponseWriter, req bunrouter.Request) error

	GetUsers(w http.ResponseWriter, req bunrouter.Request) error
	CreateUser(w http.ResponseWriter, req bunrouter.Request) error
	UpdateUser(w http.ResponseWriter, req bunrouter.Request) error
	DeleteUser(w http.ResponseWriter, req bunrouter.Request) error
	AssignUser(w http.ResponseWriter, req bunrouter.Request) error
	GetUserPermissions(w http.ResponseWriter, req bunrouter.Request) error

	GetRoles(w http.ResponseWriter, req bunrouter.Request) error
	CreateRole(w http.ResponseWriter, req bunrouter.Request) error

	GetServices(w http.ResponseWriter, req bunrouter.Request) error
	GetServicePermissions(w http.ResponseWriter, req bunrouter.Request) error
	CreateServicePermissions(w http.ResponseWriter, req bunrouter.Request) error
}

type userClient struct {
	Client
	client proto.UserServiceClient
}

func NewClient(c *core.Config) Client {
	conn, err := grpc.NewClient(*c.Url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return &userClient{client: proto.NewUserServiceClient(conn)}
}

func (sc *userClient) Login(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.LoginResponse, error) {
		data := httputils.GetBody[proto.LoginRequest](w, req)

		return sc.client.Login(req.Context(), data)
	})
}

func (sc *userClient) Validate(w http.ResponseWriter, req bunrouter.Request) error {
	authHeader := req.Header.Get("Authorization")
	token := strings.Split(authHeader, "Bearer ")[1]

	return httputils.Response(w, func() (*proto.ValidateResponse, error) {
		return sc.client.Validate(req.Context(), &proto.ValidateRequest{
			Token: token,
		})
	})
}

func (sc *userClient) Register(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.RegisterResponse, error) {
		data := httputils.GetBody[proto.RegisterRequest](w, req)

		return sc.client.Register(req.Context(), data)
	})
}

func (sc *userClient) GetUsers(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response[proto.GetUsersResponse](w, func() (*proto.GetUsersResponse, error) {
		return sc.client.GetUsers(req.Context(), &emptypb.Empty{})
	})
}

func (sc *userClient) GetRoles(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetRolesResponse, error) {
		return sc.client.GetRoles(req.Context(), &emptypb.Empty{})
	})
}

func (sc *userClient) GetServices(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetServicesResponse, error) {
		return sc.client.GetServices(req.Context(), &emptypb.Empty{})
	})
}

func (sc *userClient) CreateUser(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.CreateUserResponse, error) {
		data := httputils.GetBody[proto.CreateUserRequest](w, req)

		return sc.client.CreateUser(req.Context(), data)
	})
}
func (sc *userClient) UpdateUser(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.UpdateUserResponse, error) {
		data := httputils.GetBody[proto.UpdateUserRequest](w, req)

		return sc.client.UpdateUser(req.Context(), data)
	})
}

func (sc *userClient) DeleteUser(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.DeleteUserResponse, error) {
		data := httputils.GetBody[proto.DeleteUserRequest](w, req)
		return sc.client.DeleteUser(req.Context(), data)
	})
}

func (sc *userClient) CreateRole(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.CreateRoleResponse, error) {
		data := httputils.GetBody[proto.CreateRoleRequest](w, req)
		return sc.client.CreateRole(req.Context(), data)
	})
}

func (sc *userClient) AssignUser(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.AssignUserResponse, error) {
		data := httputils.GetBody[proto.AssignUserRequest](w, req)

		return sc.client.AssignUser(req.Context(), data)
	})
}

func (sc *userClient) CreatePermission(w http.ResponseWriter, req bunrouter.Request) error {

	return httputils.Response(w, func() (*proto.CreateServicePermissionsResponse, error) {
		data := httputils.GetBody[proto.CreateServicePermissionsRequest](w, req)

		return sc.client.CreateServicePermissions(req.Context(), data)
	})
}

func (sc *userClient) GetUserPermissions(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetUserPermissionsResponse, error) {
		data := httputils.GetBody[proto.GetUserPermissionsRequest](w, req)

		return sc.client.GetUserPermissions(req.Context(), data)
	})
}

func (sc *userClient) GetServicePermissions(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetServicePermissionsResponse, error) {
		data := httputils.GetParams[proto.GetServicePermissionsRequest](w, req)

		return sc.client.GetServicePermissions(req.Context(), data)
	})
}
