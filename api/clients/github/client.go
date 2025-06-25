package github

import (
	"fmt"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/proto"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/uptrace/bunrouter"
	"google.golang.org/grpc"
	"net/http"
)

type Client interface {
	GetCommits(w http.ResponseWriter, req bunrouter.Request) error
	GetRepositories(w http.ResponseWriter, req bunrouter.Request) error
	GetSecretContent(w http.ResponseWriter, req bunrouter.Request) error
	SyncEnvironment(w http.ResponseWriter, req bunrouter.Request) error
	DeleteSecret(w http.ResponseWriter, req bunrouter.Request) error
	CreateSecret(w http.ResponseWriter, req bunrouter.Request) error
	GetSecrets(w http.ResponseWriter, req bunrouter.Request) error
	GetPackageTags(w http.ResponseWriter, req bunrouter.Request) error
	DeletePackageVersion(w http.ResponseWriter, req bunrouter.Request) error
	GetPackages(w http.ResponseWriter, req bunrouter.Request) error
	GetPackage(w http.ResponseWriter, req bunrouter.Request) error
	GetPackageFile(w http.ResponseWriter, req bunrouter.Request) error
	CreatePackageVersion(w http.ResponseWriter, req bunrouter.Request) error
	PushPackageVersion(w http.ResponseWriter, req bunrouter.Request) error
	CreatePackage(w http.ResponseWriter, req bunrouter.Request) error
	DeletePackage(w http.ResponseWriter, req bunrouter.Request) error
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

func (svc *gitClient) GetCommits(w http.ResponseWriter, req bunrouter.Request) error {
	return GetCommitsHandler(w, req, svc.client)
}

func (svc *gitClient) GetRepositories(w http.ResponseWriter, req bunrouter.Request) error {
	return GetRepositoriesHandler(w, req, svc.client)
}

func (svc *gitClient) DeletePackage(w http.ResponseWriter, req bunrouter.Request) error {
	return DeletePackageHandler(w, req, svc.client)
}

func (svc *gitClient) GetSecretContent(w http.ResponseWriter, req bunrouter.Request) error {
	return GetSecretContentHandler(w, req, svc.client)
}

func (svc *gitClient) SyncEnvironment(w http.ResponseWriter, req bunrouter.Request) error {
	return SyncEnvironmentHandler(w, req, svc.client)
}

func (svc *gitClient) DeleteSecret(w http.ResponseWriter, req bunrouter.Request) error {
	return DeleteSecretHandler(w, req, svc.client)
}

func (svc *gitClient) CreateSecret(w http.ResponseWriter, req bunrouter.Request) error {
	return CreateSecretHandler(w, req, svc.client)
}

func (svc *gitClient) GetSecrets(w http.ResponseWriter, req bunrouter.Request) error {
	return GetSecretsHandler(w, req, svc.client)
}

func (svc *gitClient) GetPackageTags(w http.ResponseWriter, req bunrouter.Request) error {
	return GetPackageTagsHandler(w, req, svc.client)
}

func (svc *gitClient) DeletePackageVersion(w http.ResponseWriter, req bunrouter.Request) error {
	return DeletePackageVersionHandler(w, req, svc.client)
}

func (svc *gitClient) GetPackages(w http.ResponseWriter, req bunrouter.Request) error {
	return GetPackagesHandler(w, req, svc.client)
}

func (svc *gitClient) GetPackage(w http.ResponseWriter, req bunrouter.Request) error {
	return GetPackageHandler(w, req, svc.client)
}

func (svc *gitClient) GetPackageFile(w http.ResponseWriter, req bunrouter.Request) error {
	return GetPackageFileHandler(w, req, svc.client)
}

func (svc *gitClient) CreatePackageVersion(w http.ResponseWriter, req bunrouter.Request) error {
	return CreatePackageVersionHandler(w, req, svc.client)
}

func (svc *gitClient) PushPackageVersion(w http.ResponseWriter, req bunrouter.Request) error {
	return PushPackageVersionHandler(w, req, svc.client)
}

func (svc *gitClient) CreatePackage(w http.ResponseWriter, req bunrouter.Request) error {
	return CreatePackageHandler(w, req, svc.client)
}
