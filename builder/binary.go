package builder

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/heeser-io/universe-cli/shell"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/rs/zerolog"
)

type BinaryBuilder struct {
	language Language
	function v1.Function
	path     string
}

func NewBinaryBuilder(language Language, function v1.Function, path string) *BinaryBuilder {
	return &BinaryBuilder{language, function, path}
}

func (b *BinaryBuilder) Build() (bool, error) {

	functionPath := path.Join(b.path, b.function.Handler, "main.go")
	outputFolder := path.Join(b.path, "bin")
	outputPath := path.Join(outputFolder, b.function.Handler)

	shell.Call(fmt.Sprintf("rm %s/*", outputFolder))

	logger := zerolog.New(os.Stdout).With().Int64("time", time.Now().Unix()).Str("step", "build").Logger()

	loggerCtx := logger.WithContext(context.Background())

	err := b.language.build(loggerCtx, functionPath, outputPath)
	if err != nil {
		return false, err
	}

	// make zip
	if err := MakeZip(b.language, b.path, b.function); err != nil {
		return false, err
	}
	return true, nil
}
