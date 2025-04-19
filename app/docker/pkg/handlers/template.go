package handlers

import (
	"bytes"
	"embed"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/types"
	st "github.com/alpha-omega-corp/cloud/core/types"

	"html/template"
	"io/fs"
	"sync"
)

var (
	//go:embed templates
	embedFS      embed.FS
	unwrapFSOnce sync.Once
	unwrappedFS  fs.FS
)

type TemplateHandler interface {
	CreateDockerFile(content []byte) (*bytes.Buffer, error)
	CreateDockerBuild(name string, tag string) (*bytes.Buffer, error)
}

type templateHandler struct {
	TemplateHandler
	template *template.Template
	config   *st.Config
}

func NewTemplateHandler(config *st.Config) TemplateHandler {
	fileSys := getFS()
	tmpl, err := template.ParseFS(fileSys, "*.template")
	if err != nil {
		panic(err)
	}

	return &templateHandler{
		template: tmpl,
		config:   config,
	}
}

func (h *templateHandler) CreateDockerFile(content []byte) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	if err := h.template.ExecuteTemplate(buf, "dockerfile.template", types.CreateDockerFileDto{
		Content: string(bytes.Trim(content, "\x00")),
	}); err != nil {
		return nil, err
	}

	return buf, nil
}

func (h *templateHandler) CreateDockerBuild(name string, tag string) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	if err := h.template.ExecuteTemplate(buf, "makefile.template", &types.CreateDockerBuildDto{
		Name: name,
		Tag:  tag,
		Org:  h.config.Env.GetString("organization"),
	}); err != nil {
		return nil, err
	}

	return buf, nil
}

func getFS() fs.FS {
	unwrapFSOnce.Do(func() {
		fileSys, err := fs.Sub(embedFS, "templates")
		if err != nil {
			panic(err)
		}
		unwrappedFS = fileSys
	})
	return unwrappedFS
}
