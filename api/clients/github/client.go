package github

import (
	"bytes"
	"fmt"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/proto"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/alpha-omega-corp/cloud/core/httputils"
	"github.com/uptrace/bunrouter"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"mime/multipart"
	"net/http"
)

type Client interface {
	GetCommits(w http.ResponseWriter, req bunrouter.Request) error
	GetRepositories(w http.ResponseWriter, req bunrouter.Request) error

	GetSecrets(w http.ResponseWriter, req bunrouter.Request) error
	GetSecret(w http.ResponseWriter, req bunrouter.Request) error
	CreateSecret(w http.ResponseWriter, req bunrouter.Request) error
	DeleteSecret(w http.ResponseWriter, req bunrouter.Request) error
	SyncSecrets(w http.ResponseWriter, req bunrouter.Request) error

	GetPackages(w http.ResponseWriter, req bunrouter.Request) error
	GetPackage(w http.ResponseWriter, req bunrouter.Request) error
	CreatePackage(w http.ResponseWriter, req bunrouter.Request) error
	DeletePackage(w http.ResponseWriter, req bunrouter.Request) error
	GetPackageFile(w http.ResponseWriter, req bunrouter.Request) error
	GetPackageVersions(w http.ResponseWriter, req bunrouter.Request) error
	CreatePackageVersion(w http.ResponseWriter, req bunrouter.Request) error
	DeletePackageVersion(w http.ResponseWriter, req bunrouter.Request) error
}

type gitClient struct {
	Client

	client proto.GithubServiceClient
}

func NewClient(c *core.Config) Client {
	conn, err := grpc.NewClient(*c.Url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return &gitClient{client: proto.NewGithubServiceClient(conn)}
}

func (sc *gitClient) GetCommits(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetCommitsResponse, error) {
		data := httputils.GetBody[proto.GetCommitsRequest](w, req)

		return sc.client.GetCommits(req.Context(), data)
	})
}

func (sc *gitClient) GetRepositories(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetRepositoriesResponse, error) {
		data := httputils.GetBody[proto.GetRepositoriesRequest](w, req)

		return sc.client.GetRepositories(req.Context(), data)
	})
}

func (sc *gitClient) GetSecrets(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetSecretsResponse, error) {
		data := httputils.GetBody[proto.GetSecretsRequest](w, req)

		return sc.client.GetSecrets(req.Context(), data)
	})
}

func (sc *gitClient) GetSecret(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetSecretResponse, error) {
		data := httputils.GetBody[proto.GetSecretRequest](w, req)

		return sc.client.GetSecret(req.Context(), data)
	})
}

func (sc *gitClient) CreateSecret(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.CreateSecretResponse, error) {
		data := httputils.GetBody[proto.CreateSecretRequest](w, req)

		return sc.client.CreateSecret(req.Context(), data)
	})
}

func (sc *gitClient) DeleteSecret(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.DeleteSecretResponse, error) {
		data := httputils.GetBody[proto.DeleteSecretRequest](w, req)

		return sc.client.DeleteSecret(req.Context(), data)
	})
}

func (sc *gitClient) SyncSecrets(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.SyncSecretsResponse, error) {
		return sc.client.SyncSecrets(req.Context(), &emptypb.Empty{})
	})
}

func (sc *gitClient) GetPackages(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetPackagesResponse, error) {
		data := httputils.GetBody[proto.GetPackagesRequest](w, req)

		return sc.client.GetPackages(req.Context(), data)
	})
}

func (sc *gitClient) GetPackage(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetPackageResponse, error) {
		data := httputils.GetParams[proto.GetPackageRequest](w, req)

		return sc.client.GetPackage(req.Context(), data)
	})
}

func (sc *gitClient) CreatePackage(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.CreatePackageResponse, error) {
		data := httputils.GetBody[proto.CreatePackageRequest](w, req)

		return sc.client.CreatePackage(req.Context(), data)
	})
}

func (sc *gitClient) DeletePackage(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.DeletePackageResponse, error) {
		data := httputils.GetBody[proto.DeletePackageRequest](w, req)

		return sc.client.DeletePackage(req.Context(), data)
	})
}

func (sc *gitClient) GetPackageVersions(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetPackageVersionsResponse, error) {
		data := httputils.GetParams[proto.GetPackageVersionsRequest](w, req)

		return sc.client.GetPackageVersions(req.Context(), data)
	})
}

func (sc *gitClient) GetPackageFile(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.GetPackageFileResponse, error) {
		data := httputils.GetParams[proto.GetPackageFileRequest](w, req)

		return sc.client.GetPackageFile(req.Context(), data)
	})
}

func (sc *gitClient) CreatePackageVersion(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.CreatePackageVersionResponse, error) {
		contents, handler, err := req.FormFile("content")

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(contents)

		file, err := handler.Open()
		if err != nil {
			return nil, err
		}

		fileBuffer := bytes.NewBuffer(make([]byte, handler.Size))
		if _, err := io.Copy(fileBuffer, file); err != nil {
			return nil, err
		}

		return sc.client.CreatePackageVersion(req.Context(), &proto.CreatePackageVersionRequest{
			Name:    req.Params().ByName("name"),
			Tag:     req.FormValue("tag"),
			Content: fileBuffer.Bytes(),
		})
	})
}

func (sc *gitClient) DeletePackageVersion(w http.ResponseWriter, req bunrouter.Request) error {
	return httputils.Response(w, func() (*proto.DeletePackageVersionResponse, error) {
		data := httputils.GetParams[proto.DeletePackageVersionRequest](w, req)

		return sc.client.DeletePackageVersion(req.Context(), data)
	})
}
