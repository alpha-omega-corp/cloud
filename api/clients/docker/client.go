package docker

import (
	"fmt"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/proto"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/alpha-omega-corp/cloud/core/httputils"
	"github.com/uptrace/bunrouter"
	"google.golang.org/grpc"
	"mime/multipart"
	"net/http"
)

type Client interface {
	CreateUserMachine(w http.ResponseWriter, req bunrouter.Request) error
	GetContainers(w http.ResponseWriter, req bunrouter.Request) error
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

func NewClient(c *core.Config) Client {
	conn, err := grpc.NewClient(*c.Url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return &dockerClient{client: proto.NewDockerServiceClient(conn)}
}

func (sc *dockerClient) CreateUserMachine(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.CreateUserMachineResponse, error) {
		data := httputils.GetBody[proto.CreateUserMachineRequest](w, req)

		return sc.client.CreateUserMachine(req.Context(), data)
	})
}
func (sc *dockerClient) GetUserMachines(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetUserMachinesResponse, error) {
		data := httputils.GetParams[proto.GetUserMachinesRequest](w, req)

		return sc.client.GetUserMachines(req.Context(), data)
	})
}
func (sc *dockerClient) GetContainers(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetContainersResponse, error) {
		data := httputils.GetParams[proto.GetContainersRequest](w, req)

		return sc.client.GetContainers(req.Context(), data)
	})
}
func (sc *dockerClient) CreateContainer(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.CreateContainerResponse, error) {
		data := httputils.GetBody[proto.CreateContainerRequest](w, req)

		return sc.client.CreateContainer(req.Context(), data)
	})
}
func (sc *dockerClient) StartContainer(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.StartContainerResponse, error) {
		data := httputils.GetBody[proto.StartContainerRequest](w, req)

		return sc.client.StartContainer(req.Context(), data)
	})
}
func (sc *dockerClient) StopContainer(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.StopContainerResponse, error) {
		data := httputils.GetBody[proto.StopContainerRequest](w, req)

		return sc.client.StopContainer(req.Context(), data)
	})
}

func (sc *dockerClient) GetContainerLogs(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetContainerLogsResponse, error) {
		data := httputils.GetParams[proto.GetContainerLogsRequest](w, req)

		return sc.client.GetContainerLogs(req.Context(), data)
	})
}

func (sc *dockerClient) DeleteContainer(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.DeleteContainerResponse, error) {
		data := httputils.GetBody[proto.DeleteContainerRequest](w, req)

		return sc.client.DeleteContainer(req.Context(), data)
	})
}

func (sc *dockerClient) GetImage(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetImageResponse, error) {
		data := httputils.GetBody[proto.GetImageRequest](w, req)

		return sc.client.GetImage(req.Context(), data)
	})
}

func (sc *dockerClient) StoreImage(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.StoreImageResponse, error) {
		data := httputils.GetFormData[struct {
			Name    string
			Content *multipart.FileHeader
		}](w, req)

		return sc.client.StoreImage(req.Context(), &proto.StoreImageRequest{
			Name:    data.Name,
			Content: make([]byte, 0),
		})
	})
}

func (sc *dockerClient) BuildImage(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.BuildImageResponse, error) {
		data := httputils.GetBody[proto.BuildImageRequest](w, req)

		return sc.client.BuildImage(req.Context(), data)

	})
}
