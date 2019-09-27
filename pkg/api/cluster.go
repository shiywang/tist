package api

type ClusterIgniter interface {
	Start() string
	Stop() string
	Shutdown() string
	Kill(containerName string) string
}
