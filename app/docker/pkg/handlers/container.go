package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/alpha-omega-corp/cloud/core/types"
	docker "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	_ "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	_ "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"io"
	"strings"
)

type ContainerHandler interface {
	CreateFrom(ctx context.Context, path string, name string) error
	GetAll(ctx context.Context) ([]docker.Container, error)
	GetAllFrom(ctx context.Context, path string) ([]docker.Container, error)
	Start(ctx context.Context, cId string) error
	Stop(ctx context.Context, cId string) error
	Delete(ctx context.Context, cId string) error
	GetLogs(containerId string, ctx context.Context) (io.ReadCloser, error)
}

type containerHandler struct {
	ContainerHandler

	client *client.Client
	config types.Config
}

func NewContainerHandler(cli *client.Client, c types.Config) ContainerHandler {
	return &containerHandler{
		client: cli,
		config: c,
	}
}

func (h *containerHandler) Start(ctx context.Context, cId string) error {
	return h.client.ContainerStart(ctx, cId, container.StartOptions{})
}

func (h *containerHandler) Stop(ctx context.Context, cId string) error {
	return h.client.ContainerStop(ctx, cId, container.StopOptions{})
}

func (h *containerHandler) Delete(ctx context.Context, cId string) error {
	return h.client.ContainerRemove(ctx, cId, container.RemoveOptions{
		Force: true,
	})
}

func (h *containerHandler) GetAll(ctx context.Context) ([]docker.Container, error) {
	containers, err := h.client.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (h *containerHandler) GetLogs(containerId string, ctx context.Context) (io.ReadCloser, error) {
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

func (h *containerHandler) PullImage(imgName string, ctx context.Context) error {
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

func (h *containerHandler) GetAllFrom(ctx context.Context, path string) ([]docker.Container, error) {
	filter := filters.NewArgs(filters.KeyValuePair{Key: "ancestor", Value: h.imageName(path)})
	return h.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filter,
	})
}

func (h *containerHandler) CreateFrom(ctx context.Context, path string, name string) error {
	imgName := h.imageName(path)
	if err := h.PullImage(imgName, ctx); err != nil {
		return err
	}

	resp, err := h.client.ContainerCreate(ctx, &container.Config{
		Image: imgName,
	}, nil, nil, nil, name)
	if err != nil {
		panic(err)
	}

	if err := h.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	return nil
}

func (h *containerHandler) imageName(path string) string {
	return h.config.Env.GetString("registry") + "/" + h.config.Env.GetString("name") + "/" + strings.Replace(path, "/", ":", 1)
}
