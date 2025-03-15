package types

type CreateDockerFileDto struct {
	Content string
}

type CreateDockerBuildDto struct {
	Name string
	Tag  string
	Org  string
}
