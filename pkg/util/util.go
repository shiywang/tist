package util

import (
	"errors"
	"fmt"
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
