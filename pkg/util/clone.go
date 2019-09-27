package util

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

var b = &BashExec{}

func GitCloneDockerCompose() string {
	green := color.FgGreen.Render
	dir, err := os.Getwd()
	if err != nil {
		CheckErr(err)
	}

	if _, err := os.Stat(dir + dockerComposeRelativePath); os.IsNotExist(err) {
		fmt.Print(green("git clone docker-compose repo....."))
		_, err = b.Exec("git", "clone", dockerComposeRepo, dir+dockerComposeRelativePath)
		CheckErr(err)
		fmt.Println(green("done"))
	}

	return dir + dockerComposeRelativePath + dockerComposeFile
}
