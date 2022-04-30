package builder

import (
	"log"
	"os"
	"time"

	"github.com/fatih/color"

	apigw "github.com/heeser-io/api-gateway/api-go/v1"
	gateway "github.com/heeser-io/api-gateway/backend/functions/gateway"
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

	// load cache

	cache := LoadOrCreate(stack.Project)
	if cache.Functions == nil {
		cache.Functions = make(map[string]*FunctionCache)
	}

	if cache.Gateways == nil {
		cache.Gateways = make(map[string]*GatewayCache)
	}

	for _, function := range stack.Functions {
		checksum := Checksum(function.Handler)
		// functions
		cf := cache.Functions[function.Name]
		if cf != nil {

			if checksum != cf.Checksum {
				log.Println("other")

				cf.Checksum = checksum
			} else {
				color.Yellow("function %s file not changed", function.Name)
			}
		} else {
			// Create function
			functionObj, err := CreateFunction(&CreateAndUploaFunction{
				Filepath: function.Handler,
				Name:     function.Name,
				Language: "golang",
			})
			if err != nil {
				color.Red("cannot create function %s, reason: %s", function.Name, err.Error())
				panic(err)
			}

			color.Green("successfully created function %s (%s)", functionObj.Name, functionObj.ID)
			cf = &FunctionCache{
				ID:       functionObj.ID,
				Checksum: checksum,
				Name:     function.Name,
			}
			cache.Functions[function.Name] = cf

			cache.LastUploaded = time.Now().Format(time.RFC3339)
		}
	}

	apiKey := os.Getenv("API_KEY")
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
			gwClient := apigw.WithAPIKey(apiKey)
			createGatewayParams := apigw.CreateGatewayParams{
				Name:   gw.Name,
				Routes: routes,
			}

			gatewayObj, err := gwClient.Gateway.Create(&createGatewayParams)
			if err != nil {
				panic(err)
			}
			color.Green("successfully created gateway %s (%s)", gw.Name, gatewayObj.ID)

			cg = &GatewayCache{
				ID:    gatewayObj.ID,
				Name:  gatewayObj.Name,
				Short: gatewayObj.Short,
			}

			cache.Gateways[gatewayObj.Name] = cg
		}
	}

	// save cache after everything is done
	cache.Save()
}
