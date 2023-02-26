package builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/heeser-io/universe-cli/shell"
	"github.com/rs/zerolog"
)

type Language interface {
	isAvailable() (bool, error)
	build(ctx context.Context, filepath string, output string) error
	checksum(filepath string) (string, error)
	name() string
}

type Nodejs struct {
	Version string
}

type Golang struct {
	Version string
}

func NewLanguage(langstr string) Language {
	if langstr == "nodejs14" {
		return &Nodejs{
			Version: "14",
		}
	}
	return &Golang{
		Version: "1.19.3",
	}
}

func (nodejs *Nodejs) isAvailable() (bool, error) {
	return true, nil
}
func (nodejs *Nodejs) version() string {
	return nodejs.Version
}
func (nodejs *Nodejs) name() string {
	return "nodejs"
}
func (nodejs *Nodejs) checksum(filepath string) (string, error) {
	return "", nil
}

func (nodejs *Nodejs) build(ctx context.Context, filepath, output string) error {
	return nil
}

func (golang *Golang) version() string {
	return golang.Version
}

func (golang *Golang) name() string {
	return "golang"
}
func (golang *Golang) checksum(filepath string) (string, error) {
	cs, err := Checksum(filepath)
	if err != nil {
		return "", err
	}
	return cs, nil
}

func (golang *Golang) isAvailable() (bool, error) {
	b, err := shell.Call("go version")
	if err != nil {
		return false, err
	}

	s := string(b)

	if strings.Contains(s, golang.Version) {
		return true, nil
	}
	return false, fmt.Errorf("golang with version %s not available", golang.Version)
}

func (golang *Golang) build(ctx context.Context, filepath string, output string) error {
	logger := zerolog.Ctx(ctx).With().Str("substep", "compile").Str("language", golang.name()).Str("version", golang.version()).Logger()

	var cmd string
	if golang.Version == "1.19.3" {
		// build filepath with go 1.18.2
		cmd = fmt.Sprintf("CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o %s %s", output, filepath)
	}
	logger.Debug().Msgf("compiling with command %s", cmd)
	_, err := shell.Call(cmd)
	if err != nil {
		logger.Err(err).Msg("")
		return err
	}

	logger.Debug().Msgf("successfully compiled %s and saved binary to: %s", filepath, output)

	// run cmd
	return nil
}
