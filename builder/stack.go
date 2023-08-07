package builder

import (
	"fmt"
	"os"
	"path"

	"github.com/heeser-io/universe-cli/config"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

type Stack struct {
	BasePath      string                `yaml:"basePath"`
	Stacks        []ServiceStack        `yaml:"stacks"`
	Project       v1.Project            `yaml:"project"`
	AuthProvider  string                `yaml:"authProvider"`
	OAuth         v1.OAuth              `yaml:"oauth"`
	Functions     []v1.Function         `yaml:"functions"`
	Gateways      []v1.Gateway          `yaml:"gateways"`
	Secrets       []v1.Secret           `yaml:"secrets"`
	Filemapping   []Filemapping         `yaml:"filemapping"`
	Tasks         []v1.CreateTaskParams `yaml:"tasks"`
	Subscriptions []v1.Subscription     `yaml:"subscriptions"`
	Templates     []v1.Template         `yaml:"templates"`
	Webhooks      []v1.Webhook          `yaml:"webhooks"`
	Domains       []v1.Domain           `yaml:"domains"`
	KeyValues     []v1.KeyValue         `yaml:"keyvalues"`
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
