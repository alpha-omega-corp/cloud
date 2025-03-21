package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/proto"
	"github.com/uptrace/bunrouter"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net/http"
	"strconv"
	"time"
)

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateRoleRequestBody struct {
	Name string `json:"name"`
}

type CreatePermissionsRequestBody struct {
	RoleID    int64 `json:"roleId"`
	ServiceID int64 `json:"serviceId"`
	CanRead   bool  `json:"canRead"`
	CanWrite  bool  `json:"canWrite"`
	CanManage bool  `json:"canManage"`
}

type CreateUserRequestBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequestBody struct {
	Name string `json:"name"`
}

type AssignUserRequestBody struct {
	UserId int64   `json:"userId"`
	Roles  []int64 `json:"roles"`
}

func LoginHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	data := new(LoginRequestBody)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	res, err := s.Login(req.Context(), &proto.LoginRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func RegisterHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	data := new(RegisterRequestBody)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	res, err := s.Register(req.Context(), &proto.RegisterRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func CreateRoleHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	data := new(CreateRoleRequestBody)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	res, err := s.CreateRole(req.Context(), &proto.CreateRoleRequest{
		Name: data.Name,
	})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func GetRolesHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	res, err := s.GetRoles(req.Context(), &proto.GetRolesRequest{})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func CreatePermissionHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	data := new(CreatePermissionsRequestBody)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	res, err := s.CreateServicePermissions(req.Context(), &proto.CreateServicePermissionsRequest{
		RoleId:    data.RoleID,
		ServiceId: data.ServiceID,
		CanRead:   data.CanRead,
		CanWrite:  data.CanWrite,
		CanManage: data.CanManage,
	})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func GetServices(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	res, err := s.GetServices(req.Context(), &proto.GetServicesRequest{})
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func GetServicePermissionsHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	serviceId, err := strconv.ParseInt(req.Params().ByName("serviceId"), 10, 64)
	if err != nil {
		return err
	}

	res, err := s.GetServicePermissions(req.Context(), &proto.GetServicePermissionsRequest{
		ServiceId: serviceId,
	})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func GetUserPermissionsHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	userId, err := strconv.ParseInt(req.Params().ByName("id"), 10, 64)
	if err != nil {
		return err
	}

	res, err := s.GetUserPermissions(req.Context(), &proto.GetUserPermissionsRequest{
		UserId: userId,
	})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func GetUsersHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	res, err := s.GetUsers(req.Context(), &proto.GetUsersRequest{})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func CreateUserHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	data := new(CreateUserRequestBody)

	fmt.Print(req.Body)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	res, err := s.CreateUser(req.Context(), &proto.CreateUserRequest{
		Name:  data.Name,
		Email: data.Email,
	})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func UpdateUserHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	userId, err := strconv.ParseInt(req.Params().ByName("id"), 10, 64)
	if err != nil {
		return err
	}

	data := new(UpdateUserRequestBody)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	res, err := s.UpdateUser(req.Context(), &proto.UpdateUserRequest{
		Id:   userId,
		Name: data.Name,
	})
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func DeleteUserHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	userId, err := strconv.ParseInt(req.Params().ByName("id"), 10, 64)
	if err != nil {
		return err
	}

	res, err := s.DeleteUser(req.Context(), &proto.DeleteUserRequest{Id: userId})
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func AssignUserHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	data := new(AssignUserRequestBody)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	res, err := s.AssignUser(req.Context(), &proto.AssignUserRequest{
		UserId: data.UserId,
		Roles:  data.Roles,
	})

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}

func GetTestHandler(w http.ResponseWriter, req bunrouter.Request, s proto.UserServiceClient) error {
	fmt.Println(req.Body)
	config := clientv3.Config{
		Endpoints:   []string{"etcd:2380"},
		DialTimeout: 5 * time.Second,
	}

	c, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
	}

	res, err := c.Put(context.Background(), "toto", "hello")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(res)

	r, err := c.Get(context.Background(), "toto")
	fmt.Println(r.Kvs)

	return bunrouter.JSON(w, &proto.UpdateUserRequest{
		Id:   http.StatusOK,
		Name: r.Kvs[0].String(),
	})
}
