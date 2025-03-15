package pkg

import (
	"context"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/handlers"
	proto2 "github.com/alpha-omega-corp/cloud/app/user/pkg/proto"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/utils"
	"github.com/uptrace/bun"
)

type Server struct {
	proto2.UnimplementedUserServiceServer

	authService handlers.AuthService
	permService handlers.PermService
	roleService handlers.RoleService
	userService handlers.UserService
}

func NewServer(db *bun.DB, w *utils.AuthWrapper) *Server {
	return &Server{
		authService: handlers.NewAuthService(w, db),
		permService: handlers.NewPermService(db),
		roleService: handlers.NewRoleService(db),
		userService: handlers.NewUserService(db),
	}
}

func (s *Server) CreateUser(ctx context.Context, req *proto2.CreateUserRequest) (*proto2.CreateUserResponse, error) {
	return s.userService.Create(ctx, req)
}
func (s *Server) GetUsers(ctx context.Context, req *proto2.GetUsersRequest) (*proto2.GetUsersResponse, error) {
	return s.userService.GetAll(ctx)
}
func (s *Server) UpdateUser(ctx context.Context, req *proto2.UpdateUserRequest) (*proto2.UpdateUserResponse, error) {
	return s.userService.Update(ctx, req)
}
func (s *Server) DeleteUser(ctx context.Context, req *proto2.DeleteUserRequest) (*proto2.DeleteUserResponse, error) {
	return s.userService.Delete(ctx, req)
}
func (s *Server) AssignUser(ctx context.Context, req *proto2.AssignUserRequest) (*proto2.AssignUserResponse, error) {
	return s.userService.Assign(ctx, req)
}
func (s *Server) GetServices(ctx context.Context, req *proto2.GetServicesRequest) (*proto2.GetServicesResponse, error) {
	return s.permService.GetServices(ctx)
}
func (s *Server) CreateServicePermission(ctx context.Context, req *proto2.CreateServicePermissionsRequest) (*proto2.CreateServicePermissionsResponse, error) {
	return s.permService.CreateServicePermissions(ctx, req)
}
func (s *Server) GetServicePermissions(ctx context.Context, req *proto2.GetServicePermissionsRequest) (*proto2.GetServicePermissionsResponse, error) {
	return s.permService.GetServicePermissions(ctx, req)
}
func (s *Server) GetUserPermissions(ctx context.Context, req *proto2.GetUserPermissionsRequest) (*proto2.GetUserPermissionsResponse, error) {
	return s.permService.GetUserPermissions(ctx, req)
}
func (s *Server) Login(ctx context.Context, req *proto2.LoginRequest) (*proto2.LoginResponse, error) {
	return s.authService.Login(ctx, req)
}
func (s *Server) Register(ctx context.Context, req *proto2.RegisterRequest) (*proto2.RegisterResponse, error) {
	return s.authService.Register(ctx, req)
}
func (s *Server) Validate(ctx context.Context, req *proto2.ValidateRequest) (*proto2.ValidateResponse, error) {
	return s.authService.Validate(ctx, req)
}
func (s *Server) GetRoles(ctx context.Context, req *proto2.GetRolesRequest) (*proto2.GetRolesResponse, error) {
	return s.roleService.GetAll(ctx)
}
func (s *Server) CreateRole(ctx context.Context, req *proto2.CreateRoleRequest) (*proto2.CreateRoleResponse, error) {
	return s.roleService.Create(ctx, req)
}
