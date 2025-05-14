package docker

import (
	"github.com/uptrace/bunrouter"
)

func RegisterClient(client Client, r *bunrouter.Router) Client {

	r.GET("/docker/containers", client.GetContainers)
	r.POST("/docker/containers", client.CreateContainer)

	r.DELETE("/docker/containers/:id", client.DeleteContainer)
	r.GET("/docker/containers/:id/logs", client.GetContainerLogs)
	r.POST("/docker/containers/:id/start", client.StartContainer)
	r.POST("/docker/containers/:id/stop", client.StopContainer)

	r.POST("/docker/images", client.StoreImage)
	r.GET("/docker/images/:name", client.GetImage)
	r.POST("/docker/images/build", client.BuildImage)

	r.GET("/docker/pkg/:name/containers/:tag", client.GetPackageVersionContainers)

	return client
}
