package pkg

import (
	"context"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/handlers"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/proto"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/utils"
	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	proto.UnimplementedUserServiceServer

	authService handlers.AuthService
	userService handlers.UserService
}

func NewServer(db *bun.DB, w *utils.AuthWrapper) *Server {
	return &Server{
		authService: handlers.NewAuthService(w, db),
		userService: handlers.NewUserService(db),
	}
}

func (s *Server) GetUsers(ctx context.Context, _ *emptypb.Empty) (*proto.GetUsersResponse, error) {
	return s.userService.GetAll(ctx)
}

func (s *Server) GetRoles(ctx context.Context, _ *emptypb.Empty) (*proto.GetRolesResponse, error) {
	return s.authService.GetRoles(ctx)
}

func (s *Server) GetServices(ctx context.Context, _ *emptypb.Empty) (*proto.GetServicesResponse, error) {
	return s.authService.GetServices(ctx)
}

func (s *Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	return s.userService.Create(ctx, req)
}

func (s *Server) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	return s.userService.Update(ctx, req)
}
func (s *Server) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	return s.userService.Delete(ctx, req)
}
func (s *Server) AssignUser(ctx context.Context, req *proto.AssignUserRequest) (*proto.AssignUserResponse, error) {
	return s.userService.Assign(ctx, req)
}
func (s *Server) GetUserPermissions(ctx context.Context, req *proto.GetUserPermissionsRequest) (*proto.GetUserPermissionsResponse, error) {
	return s.userService.GetPermissions(ctx, req)
}

func (s *Server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return s.authService.Login(ctx, req)
}
func (s *Server) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return s.authService.Register(ctx, req)
}
func (s *Server) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	return s.authService.Validate(ctx, req)
}

func (s *Server) CreateRole(ctx context.Context, req *proto.CreateRoleRequest) (*proto.CreateRoleResponse, error) {
	return s.authService.CreateRole(ctx, req)
}
