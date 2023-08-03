package builder

import (
	"fmt"
	"os"
	"path"

	"github.com/heeser-io/universe-cli/config"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

type Stack struct {
	BasePath      string                `yaml:"basePath"`
	Stacks        []ServiceStack        `yaml:"stacks"`
	Project       v2.Project            `yaml:"project"`
	OAuth         v2.OAuth              `yaml:"oauth"`
	Functions     []v2.Function         `yaml:"functions"`
	Gateways      []v2.Gateway          `yaml:"gateways"`
	Secrets       []v2.Secret           `yaml:"secrets"`
	Filemapping   []Filemapping         `yaml:"filemapping"`
	Tasks         []v2.CreateTaskParams `yaml:"tasks"`
	Subscriptions []v2.Subscription     `yaml:"subscriptions"`
	Templates     []v2.Template         `yaml:"templates"`
	Webhooks      []v2.Webhook          `yaml:"webhooks"`
	Domains       []v2.Domain           `yaml:"domains"`
	KeyValues     []v2.KeyValue         `yaml:"keyvalues"`
}

type ServiceStack struct {
	Name     string `yaml:"name"`
	BasePath string `yaml:"basePath"`
}

func GetStackFile(p string) string {
	sf := config.Main.GetString("stackFile")

	if funk.IsZero(sf) {
		sf = "universe.yml"
	}

	return path.Join(p, sf)
}

func ReadStack(filepath string) (*Stack, error) {
	filedata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("stack file %s not found. project initialized?", filepath)
	}

	s := &Stack{}
	if err := yaml.Unmarshal(filedata, s); err != nil {
		return nil, fmt.Errorf("stack file %s seems to be no valid yaml, error: %v", filepath, err)
	}
	return s, nil
}
