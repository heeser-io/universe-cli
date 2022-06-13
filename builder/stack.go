package builder

import (
	"os"

	v1 "github.com/heeser-io/universe/api/v1"
	"gopkg.in/yaml.v2"
)

type Stack struct {
	Project     string          `yaml:"project"`
	Version     string          `yaml:"version"`
	AutoRelease bool            `yaml:"auto_release"`
	OAuth       v1.OAuth        `yaml:"oauth"`
	Collections []v1.Collection `yaml:"collections"`
	Functions   []v1.Function   `yaml:"functions"`
	Gateways    []v1.Gateway    `yaml:"gateways"`
	Secrets     []v1.Secret     `yaml:"secrets"`
	Files       []v1.File       `yaml:"files"`
	Tasks       []v1.Task       `yaml:"tasks"`
}

func ReadStack(filepath string) *Stack {
	filedata, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	s := &Stack{}
	if err := yaml.Unmarshal(filedata, s); err != nil {
		panic(err)
	}
	return s
}
