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
	Path      string
	Bash      util.Execer
	Container docker.Container
	binName   string
}

func (d *DockerCompose) isUp() bool {
	d.binName = binDocker
	out, err := d.Bash.Exec(d.binName, "ps")
	util.CheckErr(err)
	if strings.Contains(out, "pingcap") {
		return true
	}
	return false
}

func (d *DockerCompose) Start() string {
	yellow := color.FgYellow.Render

	if d.isUp() {
		fmt.Println(yellow("cluster already up..."))
		return ""
	}
	d.binName = binCompose

	out, err := d.Bash.Exec(d.binName, "-f", d.Path, "up", "-d")
	util.CheckErr(err)
	return out
}

func (d *DockerCompose) Stop() string {
	d.binName = binCompose
	out, err := d.Bash.Exec(d.binName, "-f", d.Path, "stop")
	util.CheckErr(err)
	return out
}

func (d *DockerCompose) Shutdown() string {
	d.binName = binCompose
	out, err := d.Bash.Exec(d.binName, "-f", d.Path, "down")
	util.CheckErr(err)
	return out
}

func (d *DockerCompose) Kill(containerName string) string {
	if d.Container == nil {
		d.Container = docker.CreateClient()
	}

	container, err := d.Container.GetContainerByName(containerName)
	if err != nil {
		util.CheckErr(errors.New(fmt.Sprintf("get Container named " + containerName + " fail")))
		return ""
	}

	err = d.Container.StopContainer(container.ID)
	if err != nil {
		util.CheckErr(errors.New(fmt.Sprintf("stop Container " + containerName + " fail")))
		return ""
	}
	return "kill " + container.ID + " successfully"
}
