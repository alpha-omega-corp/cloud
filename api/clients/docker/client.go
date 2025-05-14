package docker

import (
	"fmt"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/proto"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/uptrace/bunrouter"
	"google.golang.org/grpc"
	"net/http"
)

type Client interface {
	GetContainers(w http.ResponseWriter, req bunrouter.Request) error
	GetPackageVersionContainers(w http.ResponseWriter, req bunrouter.Request) error
	CreateContainer(w http.ResponseWriter, req bunrouter.Request) error
	GetContainerLogs(w http.ResponseWriter, req bunrouter.Request) error
	DeleteContainer(w http.ResponseWriter, req bunrouter.Request) error
	StartContainer(w http.ResponseWriter, req bunrouter.Request) error
	StopContainer(w http.ResponseWriter, req bunrouter.Request) error
	GetImage(w http.ResponseWriter, req bunrouter.Request) error
	StoreImage(w http.ResponseWriter, req bunrouter.Request) error
	BuildImage(w http.ResponseWriter, req bunrouter.Request) error
}

type dockerClient struct {
	Client
	client proto.DockerServiceClient
}

func NewClient(c *config.Config) Client {
	conn, err := grpc.Dial(*c.Url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return &dockerClient{client: proto.NewDockerServiceClient(conn)}
}

func (svc *dockerClient) GetContainers(w http.ResponseWriter, req bunrouter.Request) error {
	return GetContainersHandler(w, req, svc.client)
}

func (svc *dockerClient) CreateContainer(w http.ResponseWriter, req bunrouter.Request) error {
	return CreateContainerHandler(w, req, svc.client)
}

func (svc *dockerClient) StartContainer(w http.ResponseWriter, req bunrouter.Request) error {
	return StartContainerHandler(w, req, svc.client)
}

func (svc *dockerClient) StopContainer(w http.ResponseWriter, req bunrouter.Request) error {
	return StopContainerHandler(w, req, svc.client)
}

func (svc *dockerClient) GetPackageVersionContainers(w http.ResponseWriter, req bunrouter.Request) error {
	return GetContainersHandler(w, req, svc.client)
}

func (svc *dockerClient) GetContainerLogs(w http.ResponseWriter, req bunrouter.Request) error {
	return GetContainerLogsHandler(w, req, svc.client)
}

func (svc *dockerClient) DeleteContainer(w http.ResponseWriter, req bunrouter.Request) error {
	return DeleteContainerHandler(w, req, svc.client)
}

func (svc *dockerClient) GetImage(w http.ResponseWriter, req bunrouter.Request) error {
	return GetImageHandler(w, req, svc.client)
}

func (svc *dockerClient) StoreImage(w http.ResponseWriter, req bunrouter.Request) error {
	return StoreImageHandler(w, req, svc.client)
}

func (svc *dockerClient) BuildImage(w http.ResponseWriter, req bunrouter.Request) error {
	return BuildImageHandler(w, req, svc.client)
}
