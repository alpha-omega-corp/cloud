package pkg

import (
	"context"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/handlers"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/proto"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/docker/docker/client"
	"github.com/uptrace/bun"
	"net/http"
)

type Server struct {
	proto.UnimplementedDockerServiceServer
	imageService     handlers.ImageService
	containerService handlers.ContainerService
}

func NewServer(config *config.Config, client *client.Client, db *bun.DB) *Server {
	return &Server{
		imageService:     handlers.NewImageService(config, client, db),
		containerService: handlers.NewContainerHandler(config, client),
	}
}

func (s *Server) CreateContainer(ctx context.Context, req *proto.CreateContainerRequest) (*proto.CreateContainerResponse, error) {

	if err := s.containerService.CreateContainer(ctx, req.Path, req.Name); err != nil {
		return &proto.CreateContainerResponse{
			Status: http.StatusOK,
		}, err
	}

	return &proto.CreateContainerResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) GetImage(ctx context.Context, req *proto.GetImageRequest) (*proto.GetImageResponse, error) {
	return s.imageService.GetImage(ctx, req)
}

func (s *Server) StoreImage(ctx context.Context, req *proto.StoreImageRequest) (*proto.StoreImageResponse, error) {
	return s.imageService.StoreImage(ctx, req)
}

func (s *Server) BuildImage(ctx context.Context, req *proto.BuildImageRequest) (*proto.BuildImageResponse, error) {
	return s.imageService.BuildImage(ctx, req)
}

func (s *Server) GetContainers(ctx context.Context, req *proto.GetContainersRequest) (*proto.GetContainersResponse, error) {
	res, err := s.containerService.GetAllByImage(ctx, req.Path)
	if err != nil {
		return nil, err
	}

	resSlice := make([]*proto.Container, len(res))
	for index, container := range res {
		resSlice[index] = &proto.Container{
			Id:      container.ID,
			Names:   container.Names,
			Image:   container.Image,
			Status:  container.Status,
			Command: container.Command,
			State:   container.State,
			Created: container.Created,
		}
	}

	return &proto.GetContainersResponse{
		Containers: resSlice,
	}, nil

}
