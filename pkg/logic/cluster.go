package logic

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
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
	yellow := color.FgYellow.Render

	if d.isUp() {
		fmt.Println(yellow("cluster already up..."))
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
		util.CheckErr(errors.New(fmt.Sprintf("get Container named " + containerName + " fail")))
		return
	}

	err = ctl.StopContainer(container.ID)
	if err != nil {
		util.CheckErr(errors.New(fmt.Sprintf("stop container " + containerName + " fail")))
		return
	}
}
