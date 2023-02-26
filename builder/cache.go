package builder

import (
	"os"
	"path"

	"github.com/heeser-io/universe-cli/config"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

type Cache struct {
	Project       *v1.Project                 `yaml:"project"`
	Checksum      string                      `yaml:"checksum"`
	LastUploaded  string                      `yaml:"lastUploaded"`
	Secrets       map[string]*v1.Secret       `yaml:"secrets"`
	Functions     map[string]*v1.Function     `yaml:"functions"`
	Gateways      map[string]*v1.Gateway      `yaml:"gateways"`
	OAuth         *v1.OAuth                   `yaml:"oauth"`
	Filemappings  map[string][]*v1.File       `yaml:"filemappings"`
	Tasks         map[string]*v1.Task         `yaml:"tasks"`
	Subscriptions map[string]*v1.Subscription `yaml:"subscriptions"`
	Templates     map[string]*v1.Template     `yaml:"templates"`
	Webhooks      map[string]*v1.Webhook      `yaml:"webhooks"`
	Domains       map[string]*v1.Domain       `yaml:"domains"`
	path          string
}

func getCacheFile(p string) string {
	cf := config.Main.GetString("cacheFile")
	if funk.IsZero(cf) {
		cf = ".stack.yml"
	}

	return path.Join(p, cf)
}
func LoadOrCreate(p string) *Cache {

	fb, err := os.ReadFile(getCacheFile(p))

	create := false
	if err != nil {
		create = true
	}

	c := Cache{}
	if !create {
		if err := yaml.Unmarshal(fb, &c); err != nil {
			panic(err)
		}
	}

	c.path = p
	return &c
}

func (c *Cache) Save() {
	b, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(getCacheFile(c.path), b, os.ModePerm); err != nil {
		panic(err)
	}
}
