package logic

import (
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
)

// TODO: remove this code entirely : ) @shiywang
// Usually we don't fussing on unit test
// but it's useful to using interfaces mocking when it
// comes to larger projects.

type fakeBashExec struct{}

type fakeDockerAPI struct{}

func (ctl *fakeDockerAPI) GetContainerByName(name string) (*types.Container, error) {
	return &types.Container{Names: []string{"test"}, ID: "haha"}, nil
}

func (ctl *fakeDockerAPI) StopContainer(containerID string) error {
	return nil
}

func (ctl *fakeDockerAPI) GetAllStartingContainer() ([]types.Container, error) {
	return []types.Container{}, nil
}

func (f fakeBashExec) Exec(s string, args ...string) (string, error) {
	var str string
	for _, str = range args {
		str = " " + str
	}
	return s + " " + str, nil
}

func TestStart(t *testing.T) {
	expect := "docker-compose  -d"
	test := &DockerCompose{Bash: &fakeBashExec{}, Path: "./test", binName: binCompose}

	out := test.Start()

	if out != expect {
		t.Error(fmt.Sprintf("Start function return value is not expected (%s vs %s)", out, expect))
	}
}

func TestStop(t *testing.T) {
	expect := "docker-compose  stop"
	test := &DockerCompose{Bash: &fakeBashExec{}, Path: "./test", binName: binCompose}

	out := test.Stop()

	if out != expect {
		t.Error(fmt.Sprintf("Stop function return value is not expected (%s vs %s)", out, expect))
	}
}

func TestShutdown(t *testing.T) {
	expect := "docker-compose  down"
	test := &DockerCompose{Bash: &fakeBashExec{}, Path: "./test", binName: binCompose}

	out := test.Shutdown()

	if out != expect {
		t.Error(fmt.Sprintf("Shudown function return value is not expected (%s vs %s)", out, expect))
	}
}

func TestKill(t *testing.T) {
	expect := "kill haha successfully"
	test := &DockerCompose{Bash: &fakeBashExec{}, Path: "./test", binName: binCompose, Container: &fakeDockerAPI{}}

	out := test.Kill("test")

	if out != expect {
		t.Error(fmt.Sprintf("Shudown function return value is not expected (%s vs %s)", out, expect))
	}
}
