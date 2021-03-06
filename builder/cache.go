package builder

import (
	"os"

	"github.com/heeser-io/universe-cli/config"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

type Cache struct {
	Project      *v1.Project               `yaml:"project"`
	Checksum     string                    `yaml:"checksum"`
	LastUploaded string                    `yaml:"lastUploaded"`
	Secrets      map[string]*v1.Secret     `yaml:"secrets"`
	Collections  map[string]*v1.Collection `yaml:"collections"`
	Functions    map[string]*v1.Function   `yaml:"functions"`
	Gateways     map[string]*v1.Gateway    `yaml:"gateways"`
	OAuth        *v1.OAuth                 `yaml:"oauth"`
	Files        map[string]*v1.File       `yaml:"files"`
	Tasks        map[string]*v1.Task       `yaml:"tasks"`
}

func getCacheFile() string {
	cf := config.Main.GetString("cacheFile")
	if funk.IsZero(cf) {
		cf = ".stack.yml"
	}

	return cf
}
func LoadOrCreate() *Cache {

	fb, err := os.ReadFile(getCacheFile())

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

	return &c
}

func (c *Cache) Save() {
	b, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(getCacheFile(), b, os.ModePerm); err != nil {
		panic(err)
	}
}
