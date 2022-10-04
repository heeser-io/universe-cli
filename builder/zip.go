package builder

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/heeser-io/universe-cli/shell"
	v1 "github.com/heeser-io/universe/api/v1"
)

func MakeZip(lang Language, p string, function v1.Function) error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	outputPath := path.Join(p, "bin")

	if p == "" {
		p = "."
	}

	os.MkdirAll(outputPath, os.ModePerm)

	if lang.name() == "nodejs" {
		cmd := fmt.Sprintf("sh -c 'zip -r %s.zip %s -x \"./node_modules/**\\*\" -x universe.yml'", path.Join(outputPath, function.Name), p)
		log.Println(cmd)
		shell.Call(cmd)
	} else {
		shell.Call(fmt.Sprintf("cd %s && zip -r %s.zip %s && cd %s", outputPath, function.Handler, function.Handler, currentPath))
	}

	return nil
}
