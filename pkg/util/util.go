package util

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/gookit/color"
)

const dockerComposeRepo = "https://github.com/pingcap/tidb-docker-compose.git"
const dockerComposeRelativePath = "/tidb-docker-compose"
const dockerComposeFile = "/docker-compose.yml"

func CheckErr(err error) {
	red := color.FgRed.Render
	if err != nil {
		panic(red(err.Error()))
	}
}

func Exec(s string, args ...string) (string, error) {
	cmd := exec.Command(s, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func Call(m map[string]interface{}, name string, params ...string) (result []reflect.Value, err error) {
	green := color.FgGreen.Render

	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New(fmt.Sprintf("(%s) The number of params is not matched. (%d vs %d)", name, len(params), f.Type().NumIn()))
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	fmt.Println(green("start to calling step ", name))
	result = f.Call(in)
	fmt.Println(green("succeed in calling step ", name))

	return
}

func GitCloneDockerCompose() string {
	green := color.FgGreen.Render
	dir, err := os.Getwd()
	if err != nil {
		CheckErr(err)
	}

	if _, err := os.Stat(dir + dockerComposeRelativePath); os.IsNotExist(err) {
		fmt.Print(green("git clone docker-compose repo....."))
		_, err = Exec("git", "clone", dockerComposeRepo, dir+dockerComposeRelativePath)
		CheckErr(err)
		fmt.Println(green("done"))
	}

	return dir + dockerComposeRelativePath + dockerComposeFile
}
