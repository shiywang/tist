package docker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/shiywang/tist/pkg/util"
)

var once sync.Once
var ctl *DockerCtl

type Container interface {
	GetContainerByName(name string) (*types.Container, error)
	StopContainer(containerID string) error
	GetAllStartingContainer() ([]types.Container, error)
}

type DockerCtl struct {
	cli *client.Client
}

func CreateClient() *DockerCtl {
	//Singleton mode only init once.
	once.Do(func() {
		cli, err := client.NewClient(client.DefaultDockerHost, "v1.12", nil, nil)
		if err != nil {
			util.CheckErr(err)
		}
		ctl = &DockerCtl{cli}
	})

	return ctl
}

func (ctl *DockerCtl) GetContainerByName(name string) (*types.Container, error) {
	args := filters.NewArgs()
	args.Add("name", name)

	containers, err := ctl.cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: args})
	if err != nil {
		return nil, err
	}
	if len(containers) < 1 {
		return nil, errors.New("no contain has name " + name)
	}

	return &containers[0], nil
}

func (ctl *DockerCtl) StopContainer(containerID string) error {
	timeout := time.Second * 10
	err := ctl.cli.ContainerStop(context.Background(), containerID, &timeout)
	return err
}

func (ctl *DockerCtl) GetAllStartingContainer() ([]types.Container, error) {
	return ctl.cli.ContainerList(context.Background(), types.ContainerListOptions{})
}
