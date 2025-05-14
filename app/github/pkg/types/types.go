package types

type CreateMakefileDto struct {
	Registry string
	OrgName  string
	Name     string
	Tag      string
}

type CreateDockerfileDto struct {
	Name    string
	Tag     string
	Author  string
	Content string
}

type GitPackage struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"package_type"`
	Version    int64  `json:"version_count"`
	Visibility string `json:"visibility"`
	Url        string `json:"url"`
	HtmlUrl    string `json:"html_url"`

	Owner struct {
		Id     int64  `json:"id"`
		Name   string `json:"login"`
		Type   string `json:"type"`
		NodeId string `json:"node_id"`
	} `json:"owner"`
}

type GitPackageVersions struct {
	Items []GitPackageVersion `json:"items"`
}

type GitPackageVersion struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	PackageHtmlUrl string `json:"package_html_url"`
	HtmlUrl        string `json:"html_url"`
	Created        string `json:"created_at"`
	Updated        string `json:"updated_at"`
	Metadata       struct {
		PackageType string `json:"package_type"`
		Container   struct {
			Tags []string `json:"tags"`
		} `json:"container"`
	} `json:"metadata"`
}
