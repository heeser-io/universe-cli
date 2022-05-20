package builder

import (
	"os"
	"time"

	"github.com/fatih/color"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/heeser-io/universe/services/gateway"

	"gopkg.in/yaml.v2"
)

type Stack struct {
	Project     string     `yml:"project"`
	Version     string     `yml:"version"`
	AutoRelease bool       `yml:"auto_release"`
	Functions   []Function `yml:"functions"`
	Gateways    []Gateway  `yml:"gateways"`
}
type Function struct {
	Name    string `yml:"name"`
	Handler string `yml:"handler"`
	Path    string `yml:"path"`
}

type Gateway struct {
	Name   string  `yml:"name"`
	Routes []Route `yml:"routes"`
}

type Route struct {
	Function string `yml:"function"`
	Method   string `yml:"method"`
	Path     string `yml:"path"`
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

func BuildStack(filepath string) {
	stack := ReadStack(filepath)

	// Check if we already have created a function
	// Check old filehash vs new filehash
	// create if none, update if differs, do nothing if its the same

	apiKey := os.Getenv("API_KEY")
	gwClient := v1.WithAPIKey(apiKey)
	// load cache
	cache := LoadOrCreate()
	if cache.Functions == nil {
		cache.Functions = make(map[string]*v1.Function)
	}

	if cache.Gateways == nil {
		cache.Gateways = make(map[string]*v1.Gateway)
	}

	if cache.Project.ID == "" {
		projectObj, err := gwClient.Project.Create(&v1.CreateProjectParams{
			Name: stack.Project,
		})
		if err != nil {
			panic(err)
		}
		cache.Project = *projectObj
	}
	for _, function := range stack.Functions {
		checksum := Checksum(function.Path)
		// functions
		cf := cache.Functions[function.Name]
		if cf != nil {

			if checksum != cf.Checksum {
				// update file and function
				functionObj, err := UpdateFunction(&UpdateAndUploadFunction{
					FunctionID: cf.ID,
					Filepath:   function.Path,
					Checksum:   checksum,
				})

				color.Green("successfully updated function %s (%s) to version %d", functionObj.Name, functionObj.ID, functionObj.Version)
				if err != nil {
					panic(err)
				}

				// release function
				if err := ReleaseFunction(functionObj.ID); err != nil {
					panic(err)
				}

				color.Green("successfully released function %s (%s) with version %d", functionObj.Name, functionObj.ID, functionObj.Version)

				now := time.Now().Format(time.RFC3339)
				cf.LastReleasedAt = now
				cf.Checksum = checksum
				cache.LastUploaded = now
			} else {
				color.Yellow("function %s file not changed", function.Name)
			}
		} else {
			// Create function
			functionObj, err := CreateFunction(&CreateAndUploadFunction{
				Filepath:  function.Path,
				Handler:   function.Handler,
				Checksum:  checksum,
				ProjectID: cache.Project.ID,
				Name:      function.Name,
				Language:  "golang",
			})
			if err != nil {
				color.Red("cannot create function %s, reason: %s", function.Name, err.Error())
				panic(err)
			}

			color.Green("successfully created function %s (%s)", functionObj.Name, functionObj.ID)

			cf = functionObj

			// release function
			if err := ReleaseFunction(functionObj.ID); err != nil {
				panic(err)
			}
			cf.LastReleasedAt = time.Now().Format(time.RFC3339)
			cache.Functions[function.Name] = cf

			color.Green("successfully released function %s (%s) with version", functionObj.Name, functionObj.ID, functionObj.Version)
			cache.LastUploaded = time.Now().Format(time.RFC3339)
		}
	}

	for _, gw := range stack.Gateways {
		cg := cache.Gateways[gw.Name]

		if cg != nil {
			color.Yellow("gateway %s exists", gw.Name)
		} else {

			// build routes
			routes := []gateway.Route{}
			for _, r := range gw.Routes {
				routes = append(routes, gateway.Route{
					FunctionID: cache.Functions[r.Function].ID,
					Path:       r.Path,
					Method:     r.Method,
				})
			}

			createGatewayParams := v1.CreateGatewayParams{
				ProjectID: cache.Project.ID,
				Name:      gw.Name,
				Routes:    routes,
			}

			gatewayObj, err := gwClient.Gateway.Create(&createGatewayParams)
			if err != nil {
				panic(err)
			}
			color.Green("successfully created gateway %s (%s)", gw.Name, gatewayObj.ID)

			cg = gatewayObj

			cache.Gateways[gatewayObj.Name] = cg
		}
	}

	// save cache after everything is done
	cache.Save()
}
