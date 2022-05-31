package builder

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	"github.com/heeser-io/universe-cli/config"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/heeser-io/universe/services/gateway"
)

var (
	API_KEY string
)

func BuildStack() {
	stack := ReadStack(config.Main.GetString("stackFile"))

	// Check if we already have created a function
	// Check old filehash vs new filehash
	// create if none, update if differs, do nothing if its the same

	// load cache
	cache := LoadOrCreate()
	if cache.Functions == nil {
		cache.Functions = make(map[string]*v1.Function)
	}

	if cache.Gateways == nil {
		cache.Gateways = make(map[string]*v1.Gateway)
	}

	if cache.Collections == nil {
		cache.Collections = make(map[string]*v1.Collection)
	}

	if cache.Secrets == nil {
		cache.Secrets = make(map[string]*v1.Secret)
	}

	if cache.Project == nil {
		projectObj, err := client.Client.Project.Create(&v1.CreateProjectParams{
			Name: stack.Project,
		})
		if err != nil {
			panic(err)
		}
		cache.Project = projectObj
	}

	projectID := cache.Project.ID

	for _, s := range stack.Secrets {
		cg := cache.Secrets[s.Name]

		if cg != nil {
			color.Yellow("secret %s exists", s.ID)
		} else {
			createSecretParams := v1.CreateSecretParams{
				Value: s.Value,
				Tags:  s.Tags,
				Name:  s.Name,
			}

			secretObj, err := client.Client.Secret.Create(&createSecretParams)
			if err != nil {
				panic(err)
			}

			color.Green("successfully created secret %s", secretObj.ID)

			cg = secretObj
			cache.Secrets[cg.Name] = cg
		}
	}

	for _, collection := range stack.Collections {
		cc := cache.Collections[collection.Name]

		if cc != nil {

		} else {
			collectionObj, err := client.Client.Collection.Create(&v1.CreateCollectionParams{
				ProjectID: projectID,
				Name:      collection.Name,
				IndexType: collection.IndexType,
			})
			if err != nil {
				panic(err)
			}
			cc = collectionObj
			cache.Collections[collection.Name] = cc
		}
	}

	// push functions in parallel
	wg := sync.WaitGroup{}

	for _, function := range stack.Functions {
		wg.Add(1)
		go func(function v1.Function) {
			checksum := Checksum(strings.Split(function.Path, ".zip")[0])
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
					wg.Done()
				} else {
					color.Yellow("function %s file not changed", function.Name)
					wg.Done()
				}
			} else {
				// Create function
				environment := []string{}

				for _, e := range function.Environment {
					s := cache.Secrets[e]
					if s == nil {
						panic(fmt.Sprintf("secret %s not found", e))
					}
					environment = append(environment, s.ID)
				}
				functionObj, err := CreateFunction(&v1.Function{
					Path:        function.Path,
					Handler:     function.Handler,
					Checksum:    checksum,
					ProjectID:   projectID,
					Name:        function.Name,
					Tags:        function.Tags,
					Language:    function.Language,
					Environment: environment,
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
				cache.Functions[functionObj.Name] = cf

				color.Green("successfully released function %s (%s) with version %d", functionObj.Name, functionObj.ID, functionObj.Version)
				cache.LastUploaded = time.Now().Format(time.RFC3339)
				wg.Done()
			}
		}(function)
	}

	wg.Wait()
	cache.Save()

	for _, gw := range stack.Gateways {
		cg := cache.Gateways[gw.Name]

		if cg != nil {
			color.Yellow("gateway %s exists", gw.Name)

			routes := []gateway.Route{}

			for _, r := range gw.Routes {
				routes = append(routes, gateway.Route{
					FunctionID: cache.Functions[r.FunctionID].ID,
					Path:       r.Path,
					Method:     r.Method,
					AuthType:   r.AuthType,
				})
			}
			updateGatewayParams := v1.UpdateGatewayParams{
				GatewayID: cg.ID,
				Routes:    routes,
				Tags:      gw.Tags,
				Name:      gw.Name,
			}

			gatewayObj, err := client.Client.Gateway.Update(&updateGatewayParams)
			if err != nil {
				panic(err)
			}

			color.Green("successfully updated gateway %s (%s)", gw.Name, gatewayObj.ID)
			cg = gatewayObj
		} else {
			// build routes
			routes := []gateway.Route{}
			for _, r := range gw.Routes {
				// r.FunctionID
				routes = append(routes, gateway.Route{
					FunctionID: cache.Functions[r.FunctionID].ID,
					Path:       r.Path,
					Method:     r.Method,
					AuthType:   r.AuthType,
				})
			}

			createGatewayParams := v1.CreateGatewayParams{
				ProjectID: projectID,
				Name:      gw.Name,
				Routes:    routes,
				Tags:      gw.Tags,
			}

			gatewayObj, err := client.Client.Gateway.Create(&createGatewayParams)
			if err != nil {
				panic(err)
			}

			color.Green("successfully created gateway %s (%s)", gw.Name, gatewayObj.ID)

			cg = gatewayObj
		}
		cache.Gateways[gw.Name] = cg
	}

	// save cache after everything is done
	cache.Save()
}
