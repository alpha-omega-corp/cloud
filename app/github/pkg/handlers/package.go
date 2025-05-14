package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/types"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/utils"
	"os/exec"
	"path/filepath"
	"strconv"
)

type PackageService interface {
	GetVersions(name string) ([]types.GitPackageVersion, error)
	GetVersion(name string, vId int64) (*types.GitPackageVersion, error)
	Push(path string) (err error)
	Delete(name string, vId *int64) error
}

type packageService struct {
	PackageService

	client *utils.GithubApi
}

func NewPackageService(client *utils.GithubApi) PackageService {
	return &packageService{
		client: client,
	}
}

func (s *packageService) Push(path string) (err error) {
	err = s.runMakefile(path, "create")
	err = s.runMakefile(path, "tag")
	err = s.runMakefile(path, "push")

	return
}

func (s *packageService) runMakefile(path string, act string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	cmd := exec.Command("make", act)
	cmd.Dir = path

	res, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Print(string(res))

	return nil
}

func (s *packageService) GetVersions(name string) ([]types.GitPackageVersion, error) {
	res, err := s.client.Request("GET", "/packages/container/"+name+"/versions")
	if err != nil {
		return nil, err
	}

	var versions []types.GitPackageVersion
	if err := json.Unmarshal(res, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

func (s *packageService) GetVersion(name string, vId int64) (*types.GitPackageVersion, error) {
	res, err := s.client.Request("GET", "/packages/container/"+name+"/versions/"+strconv.FormatInt(vId, 10))
	if err != nil {
		return nil, err
	}

	pkg := new(types.GitPackageVersion)
	if errBuf := json.Unmarshal(res, &pkg); errBuf != nil {
		return nil, errBuf
	}

	return pkg, nil
}
