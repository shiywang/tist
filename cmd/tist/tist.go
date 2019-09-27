package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shiywang/tist/pkg/api"
	"github.com/shiywang/tist/pkg/logic"
	"github.com/shiywang/tist/pkg/util"
	"github.com/spf13/cobra"
)

type Steps []struct {
	Desc string   `yaml:"description"`
	Code string   `yaml:"code"`
	Args []string `yaml:"args"`
}

type Test struct {
	yamlPath    string
	composePath string
	steps       Steps
	dbcli       api.DB
	cluster     api.Cluster
	funcs       map[string]interface{}
}

func NewTestCommand() *Test {
	return &Test{}
}

const dockerComposePath = "/Users/shiywang/pingcap/tidb-docker-compose/docker-compose.yml"

func main() {
	o := NewTestCommand()
	rootCmd := &cobra.Command{
		Use:   "tist",
		Short: "a demo tool to 'tist' tidb",

		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete())
			util.CheckErr(o.LoadInterface())
			util.CheckErr(o.Run())
		},
	}
	rootCmd.Flags().StringVarP(&o.yamlPath, "file", "f", "", "yamlPath of the docker")
	rootCmd.Flags().StringVarP(&o.composePath, "compose", "c", dockerComposePath, "docker-compose repo path")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}

func (o *Test) Complete() error {
	if o.yamlPath == "" {
		util.CheckErr(errors.New("you must specify yaml test file"))
	}

	filename, _ := filepath.Abs(o.yamlPath)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &o.steps)
	if err != nil {
		panic(err)
	}

	return nil
}

func (o *Test) LoadInterface() error {
	//binding all the test function and struct
	o.cluster = &logic.DockerCompose{Path: o.composePath}
	o.dbcli = &logic.MysqlDB{}

	o.funcs = map[string]interface{}{
		"start-cluster": o.cluster.Start,
		"stop-cluster":  o.cluster.Stop,
		"kill-cluster":  o.cluster.Kill,
		"down-cluster":  o.cluster.Shutdown,

		"create-db": o.dbcli.CreateDB,
		"create-tb": o.dbcli.CreateTable,
		"insert-tb": o.dbcli.InsertTable,
		"query-all": o.dbcli.QueryAll,
	}

	return nil
}

func (o *Test) Run() error {
	for _, v := range o.steps {
		var err error
		if v.Args != nil {
			_, err = util.Call(o.funcs, v.Code, v.Args...)
		} else {
			_, err = util.Call(o.funcs, v.Code)
		}
		if err != nil {
			util.CheckErr(err)
		}
	}
	return nil
}
