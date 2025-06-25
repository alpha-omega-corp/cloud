package github

import (
	"github.com/uptrace/bunrouter"
)

func RegisterClient(client Client, r *bunrouter.Router) Client {

	r.GET("/github/repositories", client.GetRepositories)
	r.GET("/github/repository/:name/commits", client.GetCommits)

	r.GET("/github/secrets", client.GetSecrets)
	r.POST("/github/secrets", client.CreateSecret)
	r.GET("/github/secrets/:name", client.GetSecretContent)
	r.DELETE("/github/secrets/:name", client.DeleteSecret)
	r.POST("/github/secrets/sync", client.SyncEnvironment)

	r.GET("/github/containers", client.GetPackages)
	r.GET("/github/containers/:name", client.GetPackage)

	return client
}
