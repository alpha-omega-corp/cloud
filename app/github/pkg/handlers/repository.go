package handlers

import (
	"context"
	"fmt"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/utils"
	"github.com/google/go-github/v71/github"
)

type Content struct {
	File     *github.RepositoryContent
	Dir      []*github.RepositoryContent
	Response *github.Response
}

type PackageFile struct {
	Name    string
	SHA     string
	Content []byte
}

type RepositoryService interface {
	GetAll(ctx context.Context) ([]*github.Repository, error)
	GetCommits(ctx context.Context, repo string) ([]*github.RepositoryCommit, error)
	GetContents(ctx context.Context, repo string, path string) (content *Content, err error)
	PutContents(ctx context.Context, repo string, path string, content []byte, sha *string) error
	DeleteContents(ctx context.Context, repo string, path string, sha string) error
	GetPackageFiles(ctx context.Context, name string) ([]*PackageFile, error)
}

type repositoryService struct {
	RepositoryService
	client *utils.GithubApi
}

func NewRepositoryService(client *utils.GithubApi) RepositoryService {
	return &repositoryService{
		client: client,
	}
}

func (h *repositoryService) GetAll(ctx context.Context) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{}
	items, _, err := h.client.Default.Repositories.ListByOrg(ctx, h.client.Org, opt)

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (h *repositoryService) GetCommits(ctx context.Context, repo string) ([]*github.RepositoryCommit, error) {
	items, _, err := h.client.Default.Repositories.ListCommits(ctx, h.client.Org, "cloud", &github.CommitsListOptions{})

	fmt.Print(items)
	if err != nil {
		panic(err)
	}

	return items, nil
}

func (h *repositoryService) GetPackageFiles(ctx context.Context, name string) ([]*PackageFile, error) {
	_, dir, _, err := h.client.Default.Repositories.GetContents(ctx, h.client.Org, "container-images", name, nil)
	if err != nil {
		return nil, err
	}

	files := make([]*PackageFile, len(dir))
	for index, file := range dir {
		f, _, _, err := h.client.Default.Repositories.GetContents(ctx, h.client.Org, "container-images", name+"/"+*file.Name, nil)
		if err != nil {
			return nil, err
		}

		content, err := f.GetContent()
		if err != nil {
			return nil, err
		}

		files[index] = &PackageFile{
			SHA:     *file.SHA,
			Name:    *file.Name,
			Content: []byte(content),
		}
	}

	return files, nil
}

func (h *repositoryService) GetContents(ctx context.Context, repo string, path string) (*Content, error) {
	file, dir, res, err := h.client.Default.Repositories.GetContents(ctx, h.client.Org, repo, path, nil)

	if err != nil {
		return nil, err
	}

	return &Content{
		File:     file,
		Dir:      dir,
		Response: res,
	}, nil
}

func (h *repositoryService) PutContents(ctx context.Context, repo string, path string, content []byte, sha *string) error {
	_, _, err := h.client.Default.Repositories.CreateFile(ctx, h.client.Org, repo, path, &github.RepositoryContentFileOptions{
		Message: github.String("Webhook: Action"),
		Content: content,
		SHA:     sha,
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *repositoryService) DeleteContents(ctx context.Context, repo string, path string, sha string) error {
	_, _, err := h.client.Default.Repositories.DeleteFile(ctx, h.client.Org, repo, path, &github.RepositoryContentFileOptions{
		Message: github.String("Webhook: Action"),
		SHA:     github.String(sha),
	})

	if err != nil {
		return err
	}

	return nil
}
