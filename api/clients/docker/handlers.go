package docker

import (
	"bytes"
	"encoding/json"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/proto"
	"github.com/uptrace/bunrouter"
	"io"
	"mime/multipart"
	"net/http"
)

func StopContainerHandler(w http.ResponseWriter, req bunrouter.Request, s proto.DockerServiceClient) error {
	res, err := s.StopContainer(req.Context(), &proto.StopContainerRequest{
		ContainerId: req.Params().ByName("id"),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

func StartContainerHandler(w http.ResponseWriter, req bunrouter.Request, s proto.DockerServiceClient) error {
	res, err := s.StartContainer(req.Context(), &proto.StartContainerRequest{
		ContainerId: req.Params().ByName("id"),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

func GetContainerLogsHandler(w http.ResponseWriter, req bunrouter.Request, s proto.DockerServiceClient) error {
	res, err := s.GetContainerLogs(req.Context(), &proto.GetContainerLogsRequest{
		ContainerId: req.Params().ByName("id"),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

func DeleteContainerHandler(w http.ResponseWriter, req bunrouter.Request, s proto.DockerServiceClient) error {
	res, err := s.DeleteContainer(req.Context(), &proto.DeleteContainerRequest{
		ContainerId: req.Params().ByName("id"),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

func CreateContainerHandler(w http.ResponseWriter, req bunrouter.Request, s proto.DockerServiceClient) error {
	data := new(CreateContainerRequestBody)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	path := req.Params().ByName("name") + "/" + req.Params().ByName("tag")
	res, err := s.CreateContainer(req.Context(), &proto.CreateContainerRequest{
		Path: path,
		Name: data.ContainerName,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

func GetContainersHandler(w http.ResponseWriter, req bunrouter.Request, s proto.DockerServiceClient) error {
	res, err := s.GetContainers(req.Context(), &proto.GetContainersRequest{
		Path: req.Params().ByName("name") + "/" + req.Params().ByName("tag"),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

func GetImageHandler(w http.ResponseWriter, req bunrouter.Request, s proto.DockerServiceClient) error {
	res, err := s.GetImage(req.Context(), &proto.GetImageRequest{
		Name: req.Param("name"),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

func StoreImageHandler(w http.ResponseWriter, req bunrouter.Request, s proto.DockerServiceClient) error {
	contents, handler, err := req.FormFile("content")

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(contents)

	file, err := handler.Open()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fileBuffer := bytes.NewBuffer(make([]byte, handler.Size))
	if _, err := io.Copy(fileBuffer, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	res, err := s.StoreImage(req.Context(), &proto.StoreImageRequest{
		Name:    req.FormValue("name"),
		Content: fileBuffer.Bytes(),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

func BuildImageHandler(w http.ResponseWriter, req bunrouter.Request, s proto.DockerServiceClient) error {
	data := new(BuildImageRequest)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	res, err := s.BuildImage(req.Context(), &proto.BuildImageRequest{
		Name: data.Name,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

type CreateImageRequest struct {
	Dockerfile *multipart.FileHeader `form:"dockerfile"`
}

type BuildImageRequest struct {
	Name string `form:"name"`
}

type PushPackageRequestBody struct {
	Tag        string `json:"tag"`
	VersionSHA string `json:"sha"`
}

type DeletePackageRequestBody struct {
	Tag string `json:"tag"`
}

type CreateContainerRequestBody struct {
	ContainerName string `json:"containerName"`
}

type GetPackageVersionContainers struct {
	Path string `json:"path"`
}

type CreatePackageRequestBody struct {
	Name string `json:"name"`
}
