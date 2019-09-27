package api

type Cluster interface {
	Start()
	Stop()
	Shutdown()
	Kill(containerName string)
}
