package pkg

import (
	"context"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/handlers"
	proto2 "github.com/alpha-omega-corp/cloud/app/docker/pkg/proto"
	st "github.com/alpha-omega-corp/cloud/core/types"
	"github.com/docker/docker/client"
	"github.com/uptrace/bun"
)

type Server struct {
	proto2.UnimplementedDockerServiceServer
	imageService handlers.ImageService
}

func NewServer(config st.Config, client *client.Client, db *bun.DB) *Server {
	return &Server{
		imageService: handlers.NewImageService(config, client, db),
	}
}

func (s *Server) GetImage(ctx context.Context, req *proto2.GetImageRequest) (*proto2.GetImageResponse, error) {
	return s.imageService.GetImage(ctx, req)
}

func (s *Server) StoreImage(ctx context.Context, req *proto2.StoreImageRequest) (*proto2.StoreImageResponse, error) {
	return s.imageService.StoreImage(ctx, req)
}

func (s *Server) BuildImage(ctx context.Context, req *proto2.BuildImageRequest) (*proto2.BuildImageResponse, error) {
	return s.imageService.BuildImage(ctx, req)
}
