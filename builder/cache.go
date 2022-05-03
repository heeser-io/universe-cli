package builder

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Cache struct {
	Project      string                    `yml:"project"`
	LastUploaded string                    `yml:"lastUploaded"`
	Functions    map[string]*FunctionCache `yml:"functions"`
	Gateways     map[string]*GatewayCache  `yml:"gateways"`
}

type GatewayCache struct {
	ID     string  `yml:"id"`
	Name   string  `yml:"name"`
	Short  string  `yml:"short"`
	Routes []Route `yml:"routes"`
}
type FunctionCache struct {
	ID       string `yml:"id"`
	Name     string `yml:"name"`
	Checksum string `yml:"checksum"`
}

func LoadOrCreate(project string) *Cache {
	fb, err := os.ReadFile(fmt.Sprintf(".%s.yml", project))

	create := false
	if err != nil {
		create = true
	}

	c := Cache{
		Project: project,
	}
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
	if err := os.WriteFile(fmt.Sprintf(".%s.yml", c.Project), b, os.ModePerm); err != nil {
		panic(err)
	}
}
