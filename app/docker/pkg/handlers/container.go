package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	_ "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	_ "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/uptrace/bun"
	"io"
	"strings"
)

type ContainerService interface {
	Create(ctx context.Context, image string, name string) (*container.CreateResponse, error)
	GetOne(ctx context.Context, cId string) (container.InspectResponse, error)
	GetAll(ctx context.Context, opts container.ListOptions) ([]container.Summary, error)
	GetAllByImage(ctx context.Context, path string) ([]container.Summary, error)
	Start(ctx context.Context, cId string) error
	Stop(ctx context.Context, cId string) error
	Delete(ctx context.Context, cId string) error
	GetLogs(containerId string, ctx context.Context) (io.ReadCloser, error)
	GetTags(containerId string, ctx context.Context) ([]string, error)
}

type containerService struct {
	ContainerService

	client *client.Client
	config *core.Config
}

func NewContainerHandler(config *core.Config, client *client.Client, db *bun.DB) ContainerService {
	return &containerService{
		client: client,
		config: config,
	}
}

func (h *containerService) GetOne(ctx context.Context, cId string) (container.InspectResponse, error) {
	return h.client.ContainerInspect(ctx, cId)
}

func (h *containerService) GetAll(ctx context.Context, opts container.ListOptions) ([]container.Summary, error) {
	return h.client.ContainerList(ctx, opts)
}

func (h *containerService) Create(ctx context.Context, image string, name string) (*container.CreateResponse, error) {
	if err := h.PullImage(image, ctx); err != nil {
		return nil, err
	}

	res, err := h.client.ContainerCreate(ctx, &container.Config{
		Image: image,
	}, nil, nil, nil, name)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (h *containerService) Start(ctx context.Context, cId string) error {
	return h.client.ContainerStart(ctx, cId, container.StartOptions{})
}

func (h *containerService) Stop(ctx context.Context, cId string) error {
	return h.client.ContainerStop(ctx, cId, container.StopOptions{})
}

func (h *containerService) Delete(ctx context.Context, cId string) error {
	return h.client.ContainerRemove(ctx, cId, container.RemoveOptions{
		Force: true,
	})
}

func (h *containerService) GetAllByImage(ctx context.Context, path string) ([]container.Summary, error) {
	filter := filters.NewArgs(filters.KeyValuePair{Key: "ancestor", Value: h.imageName(path)})
	return h.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filter,
	})
}

func (h *containerService) GetLogs(containerId string, ctx context.Context) (io.ReadCloser, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		Timestamps: true,
		Details:    false,
		Since:      "40m",
	}
	logs, err := h.client.ContainerLogs(ctx, containerId, options)

	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (h *containerService) PullImage(imgName string, ctx context.Context) error {
	authConfig := registry.AuthConfig{
		Username: "packages",
		Password: h.config.Env.GetString("token"),
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}

	authString := base64.URLEncoding.EncodeToString(encodedJSON)
	_, err = h.client.ImagePull(ctx, imgName, image.PullOptions{RegistryAuth: authString})
	if err != nil {
		return err
	}

	return nil
}

func (h *containerService) imageName(path string) string {
	return h.config.Env.GetString("registry") + "/" + h.config.Env.GetString("name") + "/" + strings.Replace(path, "/", ":", 1)
}
