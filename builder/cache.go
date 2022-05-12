package builder

import (
	"os"

	v1 "github.com/heeser-io/universe/api/v1"
	"gopkg.in/yaml.v2"
)

type Cache struct {
	Project      v1.Project              `yml:"project"`
	LastUploaded string                  `yml:"lastUploaded"`
	Functions    map[string]*v1.Function `yml:"functions"`
	Gateways     map[string]*v1.Gateway  `yml:"gateways"`
}

type GatewayCache struct {
	ID        string  `yml:"id"`
	ProjectID string  `yml:"projectId"`
	Name      string  `yml:"name"`
	Short     string  `yml:"short"`
	Routes    []Route `yml:"routes"`
}
type FunctionCache struct {
	ID       string `yml:"id"`
	Name     string `yml:"name"`
	Checksum string `yml:"checksum"`
}

func LoadOrCreate() *Cache {
	fb, err := os.ReadFile(".stack.yml")

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
	if err := os.WriteFile(".stack.yml", b, os.ModePerm); err != nil {
		panic(err)
	}
}
