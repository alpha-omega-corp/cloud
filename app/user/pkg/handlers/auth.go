package handlers

import (
	"context"
	"errors"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/models"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/proto"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/utils"
	"github.com/uptrace/bun"
	"net/http"
)

type AuthService interface {
	Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error)
	Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error)
	Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error)
	GetRoles(ctx context.Context) (*proto.GetRolesResponse, error)
	GetServices(ctx context.Context) (*proto.GetServicesResponse, error)
	CreatePermissions(ctx context.Context, req *proto.CreateServicePermissionsRequest) (*proto.CreateServicePermissionsResponse, error)
	GetPermissions(ctx context.Context, req *proto.GetServicePermissionsRequest) (*proto.GetServicePermissionsResponse, error)
	CreateRole(ctx context.Context, req *proto.CreateRoleRequest) (*proto.CreateRoleResponse, error)
}

type authService struct {
	auth *utils.AuthWrapper
	db   *bun.DB
}

func NewAuthService(w *utils.AuthWrapper, db *bun.DB) AuthService {
	return &authService{
		auth: w,
		db:   db,
	}
}

func (s *authService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	_, err := s.db.NewInsert().Model(&models.User{
		Name:     req.Username,
		Email:    req.Email,
		Password: utils.HashPassword(req.Password),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &proto.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user models.User

	if err := s.db.
		NewSelect().
		Model(&user).
		Where("email = ?", req.Email).
		Scan(ctx, &user); err != nil {
		return nil, err
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		return nil, errors.New("invalid")
	}

	token, err := s.auth.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &proto.LoginResponse{
		Token: token,
		User: &proto.User{
			Id:    user.Id,
			Email: user.Email,
		},
	}, nil
}

func (s *authService) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	claims, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.db.NewSelect().Model(&user).Where("email = ?", claims.Email).Scan(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &proto.ValidateResponse{
		User: &proto.User{
			Id:    user.Id,
			Email: user.Email,
		},
	}, nil
}

func (s *authService) GetServices(ctx context.Context) (*proto.GetServicesResponse, error) {
	var services []models.Service
	if err := s.db.NewSelect().Model(&services).Scan(ctx); err != nil {
		return nil, err
	}

	var resSlice []*proto.Service
	for _, service := range services {
		resSlice = append(resSlice, &proto.Service{
			Id:   service.Id,
			Name: service.Name,
		})
	}

	return &proto.GetServicesResponse{
		Services: resSlice,
	}, nil
}

func (s *authService) CreatePermissions(ctx context.Context, req *proto.CreateServicePermissionsRequest) (*proto.CreateServicePermissionsResponse, error) {
	permissions := &models.Permission{
		Read:      req.CanRead,
		Write:     req.CanWrite,
		Manage:    req.CanManage,
		ServiceID: req.ServiceId,
		RoleId:    req.RoleId,
	}

	_, err := s.db.NewInsert().Model(permissions).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.CreateServicePermissionsResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *authService) GetPermissions(ctx context.Context, req *proto.GetServicePermissionsRequest) (*proto.GetServicePermissionsResponse, error) {
	var service models.Service
	if err := s.db.NewSelect().
		Model(&service).
		Relation("Permissions").
		Where("id = ?", req.ServiceId).
		Scan(ctx); err != nil {
		return nil, err
	}

	resSlice := make([]*proto.Permission, len(service.Permissions))
	for index, permission := range service.Permissions {
		role := new(models.Role)
		if err := s.db.NewSelect().
			Model(role).
			Where("id  = ?", permission.RoleId).
			Scan(ctx); err != nil {
			return nil, err
		}

		resSlice[index] = &proto.Permission{
			Id: permission.Id,
			Service: &proto.Service{
				Id:   service.Id,
				Name: service.Name,
			},
			Role: &proto.Role{
				Id:   role.Id,
				Name: role.Name,
			},
			CanRead:   permission.Read,
			CanWrite:  permission.Write,
			CanManage: permission.Manage,
		}
	}

	return &proto.GetServicePermissionsResponse{
		Permissions: resSlice,
	}, nil
}

func (s *authService) GetRoles(ctx context.Context) (*proto.GetRolesResponse, error) {
	var roles []*models.Role

	err := s.db.NewSelect().Model(&roles).Scan(ctx)
	if err != nil {
		return nil, err
	}

	var resSlice []*proto.Role
	for _, role := range roles {
		resSlice = append(resSlice, &proto.Role{
			Id:   role.Id,
			Name: role.Name,
		})
	}

	return &proto.GetRolesResponse{
		Roles: resSlice,
	}, nil
}

func (s *authService) CreateRole(ctx context.Context, req *proto.CreateRoleRequest) (*proto.CreateRoleResponse, error) {
	role := new(models.Role)
	role.Name = req.Name

	_, err := s.db.NewInsert().Model(role).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &proto.CreateRoleResponse{
		Status: http.StatusCreated,
	}, nil
}
