package pkg

import (
	"context"
	"encoding/json"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/handlers"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/proto"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/types"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/utils"
)

type Server struct {
	proto.UnimplementedGithubServiceServer

	packageService    handlers.PackageService
	repositoryService handlers.RepositoryService
}

func NewServer(client *utils.GithubApi) *Server {
	return &Server{
		packageService:    handlers.NewPackageService(client),
		repositoryService: handlers.NewRepositoryService(client),
	}
}

func (s *Server) GetPackages(ctx context.Context, req *proto.GetPackagesRequest) (*proto.GetPackagesResponse, error) {
	repo, err := s.repositoryService.GetContents(ctx, "container-images", ".")
	if err != nil {
		return nil, err
	}

	resSlice := make([]*proto.SimplePackage, len(repo.Dir))

	for index, pkg := range repo.Dir {
		b, err := json.Marshal(pkg)
		if err != nil {
			return nil, err
		}

		if mErr := json.Unmarshal(b, &resSlice[index]); mErr != nil {
			return nil, mErr
		}
	}

	return &proto.GetPackagesResponse{
		Packages: resSlice,
	}, nil
}

func (s *Server) GetPackage(ctx context.Context, req *proto.GetPackageRequest) (*proto.GetPackageResponse, error) {
	c, err := s.repositoryService.GetContents(ctx, "container-images", req.Name)
	versions, err := s.packageService.GetVersions(req.Name)
	if err != nil {
		return nil, err
	}

	versionMap := make(map[string]types.GitPackageVersion)

	var versionSlice []*proto.PackageVersion
	for _, version := range versions {
		for _, tag := range version.Metadata.Container.Tags {
			versionMap[tag] = version
		}
	}

	for _, dir := range c.Dir {
		if *dir.Type == "dir" {
			v := versionMap[*dir.Name]
			pkg := &proto.PackageVersion{
				RepoName:    *dir.Name,
				RepoPath:    *dir.Path,
				RepoSha:     *dir.SHA,
				RepoLink:    *dir.HTMLURL,
				VersionId:   &v.Id,
				VersionSha:  &v.Name,
				VersionLink: &v.PackageHtmlUrl,
			}

			versionSlice = append(versionSlice, pkg)
		}
	}

	return &proto.GetPackageResponse{
		Versions: versionSlice,
	}, nil
}

func (s *Server) GetPackageFile(ctx context.Context, req *proto.GetPackageFileRequest) (*proto.GetPackageFileResponse, error) {
	c, err := s.repositoryService.GetContents(ctx, "container-images", req.Path+"/"+req.Name)
	if err != nil {
		return nil, err
	}

	file, err := c.File.GetContent()
	if err != nil {
		return nil, err
	}

	return &proto.GetPackageFileResponse{
		Content: []byte(file),
	}, nil
}
