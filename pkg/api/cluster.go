package api

type ClusterIgniter interface {
	Start()
	Stop()
	Shutdown()
	Kill(containerName string)
}
