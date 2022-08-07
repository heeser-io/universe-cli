package builder

import (
	"fmt"
	"os"
	"path"

	"github.com/heeser-io/universe-cli/shell"
	v1 "github.com/heeser-io/universe/api/v1"
)

func MakeZip(p string, function v1.Function) error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	outputPath := path.Join(p, "bin")

	shell.Call(fmt.Sprintf("cd %s && zip -r %s.zip %s && cd %s", outputPath, function.Handler, function.Handler, currentPath))

	return nil
}
