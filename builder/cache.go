package builder

import (
	"os"
	"path"

	"github.com/heeser-io/universe-cli/config"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

type Cache struct {
	Project       *v2.Project                 `yaml:"project"`
	Checksum      string                      `yaml:"checksum"`
	LastUploaded  string                      `yaml:"lastUploaded"`
	Secrets       map[string]*v2.Secret       `yaml:"secrets"`
	Functions     map[string]*v2.Function     `yaml:"functions"`
	Gateways      map[string]*v2.Gateway      `yaml:"gateways"`
	OAuth         *v2.OAuth                   `yaml:"oauth"`
	Filemappings  map[string][]*v2.File       `yaml:"filemappings"`
	Tasks         map[string]*v2.Task         `yaml:"tasks"`
	Subscriptions map[string]*v2.Subscription `yaml:"subscriptions"`
	Templates     map[string]*v2.Template     `yaml:"templates"`
	Webhooks      map[string]*v2.Webhook      `yaml:"webhooks"`
	Domains       map[string]*v2.Domain       `yaml:"domains"`
	KeyValues     map[string]*v2.KeyValue     `yaml:"keyvalues"`
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
