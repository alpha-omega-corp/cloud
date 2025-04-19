package handlers

import (
	"bytes"
	"context"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/models"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/proto"
	"github.com/alpha-omega-corp/cloud/core/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/uptrace/bun"
	"net/http"
	"os"
	"os/exec"
)

type ImageService interface {
	GetImage(ctx context.Context, req *proto.GetImageRequest) (*proto.GetImageResponse, error)
	StoreImage(ctx context.Context, req *proto.StoreImageRequest) (*proto.StoreImageResponse, error)
	BuildImage(ctx context.Context, req *proto.BuildImageRequest) (*proto.BuildImageResponse, error)
}

type imageService struct {
	ImageService
	config   *types.Config
	client   *client.Client
	template TemplateHandler
	db       *bun.DB
}

func NewImageService(config *types.Config, client *client.Client, db *bun.DB) ImageService {
	return &imageService{
		db:       db,
		client:   client,
		config:   config,
		template: NewTemplateHandler(config),
	}
}

func (s *imageService) GetImage(ctx context.Context, req *proto.GetImageRequest) (*proto.GetImageResponse, error) {
	dockerfile := new(models.Dockerfile)
	err := s.db.NewSelect().Model(dockerfile).Where("name = ?", req.Name).Scan(ctx)
	if err != nil {
		return &proto.GetImageResponse{
			Status:     http.StatusNoContent,
			Dockerfile: nil,
			Image:      nil,
		}, nil
	}

	images, err := s.client.ImageList(ctx, image.ListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", s.config.Env.GetString("organization")+"/"+req.Name)),
	})
	if err != nil {
		return nil, err
	}

	dockerImage := new(proto.Image)
	if len(images) > 0 {
		item := images[0]
		dockerImage = &proto.Image{
			Id:      item.ID,
			Size:    item.Size,
			Tags:    item.RepoTags,
			Created: item.Created,
		}
	}

	return &proto.GetImageResponse{
		Status:     http.StatusOK,
		Dockerfile: dockerfile.Content,
		Image:      dockerImage,
	}, nil
}

func (s *imageService) StoreImage(ctx context.Context, req *proto.StoreImageRequest) (*proto.StoreImageResponse, error) {
	dockerfile := new(models.Dockerfile)
	err := s.db.NewSelect().Model(dockerfile).Where("name = ?", req.Name).Scan(ctx)

	content := bytes.Trim(req.Content, "\x00")

	if err != nil {
		_, err = s.db.NewInsert().Model(&models.Dockerfile{
			Name:    req.Name,
			Content: content,
		}).Exec(ctx)

		if err != nil {
			return nil, err
		}
	} else {
		_, err = s.db.
			NewUpdate().
			Model(dockerfile).
			Set("content = ?", content).
			Where("name = ?", req.Name).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &proto.StoreImageResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *imageService) BuildImage(ctx context.Context, req *proto.BuildImageRequest) (*proto.BuildImageResponse, error) {
	path := s.config.Env.GetString("build_path") + req.Name
	tag := "latest"

	makeFile, err := s.template.CreateDockerBuild(req.Name, tag)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	dockerFile := new(models.Dockerfile)
	if err := s.db.NewSelect().Model(dockerFile).Where("name = ?", req.Name).Scan(ctx); err != nil {
		return nil, err
	}

	if err = os.WriteFile(path+"/Dockerfile", dockerFile.Content, 0644); err != nil {
		return nil, err
	}

	if err = os.WriteFile(path+"/Makefile", makeFile.Bytes(), 0644); err != nil {
		return nil, err
	}

	exec.Command("make", path, "build")

	return nil, nil
}
