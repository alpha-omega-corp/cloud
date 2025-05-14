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
