package logic

import (
	"fmt"
	"strings"

	"github.com/shiywang/tist/pkg/base/docker"
	"github.com/shiywang/tist/pkg/util"
)

const binCompose = "docker-compose"
const binDocker = "docker"

type DockerCompose struct {
	Path string
}

func (d *DockerCompose) isUp() bool {
	out, err := util.Exec(binDocker, "ps")
	util.CheckErr(err)
	if strings.Contains(out, "pingcap") {
		return true
	}
	return false
}

func (d *DockerCompose) Start() {
	if d.isUp() {
		fmt.Println("cluster already up...")
		return
	}

	_, err := util.Exec(binCompose, "-f", d.Path, "up", "-d")
	util.CheckErr(err)
}

func (d *DockerCompose) Stop() {

	_, err := util.Exec(binCompose, "-f", d.Path, "stop")
	util.CheckErr(err)
}

func (d *DockerCompose) Shutdown() {
	_, err := util.Exec(binCompose, "-f", d.Path, "down")
	util.CheckErr(err)
}

func (d *DockerCompose) Kill(containerName string) {
	ctl := docker.CreateClient()
	if ctl == nil {
		return
	}

	container, err := ctl.GetContainerByName(containerName)
	if err != nil {
		fmt.Println("get Container named " + containerName + " fail")
		fmt.Println(err)
		return
	}

	err = ctl.StopContainer(container.ID)
	if err != nil {
		fmt.Println("stop container " + containerName + " fail")
		return
	}

	fmt.Println("kill container " + containerName + " ok")
}
