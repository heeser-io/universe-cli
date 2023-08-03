package builder

import (
	"fmt"
	"os"
	"path"

	"github.com/heeser-io/universe-cli/shell"
	v2 "github.com/heeser-io/universe/api/v2"
)

func MakeZip(lang Language, p string, function v2.Function) error {
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
		// we dont need to go to outputPath, because there is no compiled executable
		cmd := fmt.Sprintf("sh -c 'cd %s && zip -r %s.zip %s %s && cd %s'", path.Join(p, function.Name), path.Join("..", "bin", function.Name), "src", "package.json", currentPath)
		_, err := shell.Call(cmd)
		if err != nil {
			return err
		}
	} else {
		shell.Call(fmt.Sprintf("cd %s && zip -r %s.zip %s && cd %s", outputPath, function.Handler, function.Handler, currentPath))
	}

	return nil
}
