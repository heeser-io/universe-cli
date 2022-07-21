package builder

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/heeser-io/universe/services/gateway"
	"github.com/thoas/go-funk"
)

type Buildable interface {
	Exists() bool
	Create() error
	Update() error
	Delete() error
}

type Builder struct {
	cache *Cache
	stack *Stack
}

func New() *Builder {
	stack := ReadStack(GetStackFile())
	cache := LoadOrCreate()

	if cache.Project == nil {
		projectObj, err := client.Client.Project.Create(&v1.CreateProjectParams{
			Name: stack.Project.Name,
			Tags: stack.Project.Tags,
		})
		if err != nil {
			panic(err)
		}
		cache.Project = projectObj
		color.Green("successfully created project %s", projectObj.ID)
	}

	if cache.Functions == nil {
		cache.Functions = make(map[string]*v1.Function)
	}

	if cache.Gateways == nil {
		cache.Gateways = make(map[string]*v1.Gateway)
	}

	if cache.Functions == nil {
		cache.Functions = make(map[string]*v1.Function)
	}

	if cache.Collections == nil {
		cache.Collections = make(map[string]*v1.Collection)
	}

	if cache.Secrets == nil {
		cache.Secrets = make(map[string]*v1.Secret)
	}

	if cache.Tasks == nil {
		cache.Tasks = make(map[string]*v1.Task)
	}

	if cache.Filemappings == nil {
		cache.Filemappings = make(map[string][]*v1.File)
	}

	cache.Save()

	return &Builder{
		cache: cache,
		stack: stack,
	}
}

// returns the projectID of the project
func (b *Builder) getProjectID() string {
	if b.cache == nil || b.cache.Project == nil {
		panic("no project in current stack")
	}
	return b.cache.Project.ID
}

// buildFunctions will try to create or update all functions in the current stack
func (b *Builder) buildFunctions() error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

	// push functions in parallel
	wg := sync.WaitGroup{}

	for _, function := range stack.Functions {
		wg.Add(1)
		go func(function v1.Function) error {
			checksum := Checksum(function.Path)
			// functions
			cf := cache.Functions[function.Name]
			environment := function.Environment

			if funk.IsZero(environment) {
				environment = map[string]string{}
			}

			for k, v := range environment {
				if strings.Contains(v, "secret:") {
					s := strings.Split(v, ":")
					if len(s) == 2 {
						secretObj := cache.Secrets[s[1]]
						if secretObj != nil {
							environment[k] = fmt.Sprintf("secret:%s", secretObj.ID)
						}
					}
				} else {
					environment[k] = v
				}
			}

			if cf != nil {
				// update file and function

				if checksum != cf.Checksum {
					functionObj, err := UpdateFunction(&v1.UpdateFunctionParams{
						FunctionID:  cf.ID,
						Path:        function.Path,
						Checksum:    &checksum,
						Name:        function.Name,
						Environment: environment,
						Tags:        function.Tags,
					})
					if err != nil {
						color.Red("cannot update function %s, reason: %s", function.Name, err.Error())
						return err
					}
					color.Green("successfully updated function %s (%s)", functionObj.Name, functionObj.ID)

					// release function
					if err := ReleaseFunction(functionObj.ID); err != nil {
						color.Red("cannot release function %s, reason: %s", function.Name, err.Error())
						return err
					}

					color.Green("successfully released function %s (%s)", functionObj.Name, functionObj.ID)

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
				functionObj, err := CreateFunction(&v1.Function{
					Path:        function.Path,
					Handler:     function.Handler,
					Checksum:    checksum,
					ProjectID:   projectID,
					Permissions: function.Permissions,
					BaseImage:   function.BaseImage,
					Name:        function.Name,
					Tags:        function.Tags,
					Language:    function.Language,
					Environment: environment,
				})
				if err != nil {
					color.Red("cannot create function %s, reason: %s", function.Name, err.Error())
					return err
				}

				color.Green("successfully created function %s (%s)", functionObj.Name, functionObj.ID)

				cf = functionObj

				// release function
				if err := ReleaseFunction(functionObj.ID); err != nil {
					color.Red("cannot release function %s, reason: %s", function.Name, err.Error())
					return err
				}
				cf.LastReleasedAt = time.Now().Format(time.RFC3339)
				cache.Functions[functionObj.Name] = cf

				color.Green("successfully released function %s (%s) with version %d", functionObj.Name, functionObj.ID, functionObj.Version)
				cache.LastUploaded = time.Now().Format(time.RFC3339)
				wg.Done()
			}
			return nil
		}(function)
	}

	wg.Wait()

	cache.Save()
	return nil
}

// buildGateways will try to create or update all gateways in the current stack
func (b *Builder) buildGateways() error {
	cache := b.cache
	stack := b.stack
	projectID := b.getProjectID()

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

	cache.Save()
	return nil
}

// buildTasks will try to create or update all tasks in the current stack
func (b *Builder) buildTasks() error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

	for _, t := range stack.Tasks {
		ct := cache.Tasks[t.Name]

		if ct != nil {
			color.Yellow("task %s exists", ct.Name)

			functionID := t.FunctionID

			if strings.Contains(t.FunctionID, "function:") {
				s := strings.Split(t.FunctionID, "function:")
				if len(s) == 2 {
					functionID = s[1]
				} else {
					panic("wrong format")
				}
			}

			updateTaskParams := v1.UpdateTaskParams{
				TaskID:     ct.ID,
				FunctionID: cache.Functions[functionID].ID,
				Interval:   t.Interval,
			}
			taskObj, err := client.Client.Task.Update(&updateTaskParams)
			if err != nil {
				panic(err)
			}
			color.Green("successfully updated task %s", taskObj.ID)
			ct = taskObj
			cache.Tasks[ct.Name] = ct
		} else {
			functionID := t.FunctionID

			if strings.Contains(t.FunctionID, "function:") {
				s := strings.Split(t.FunctionID, "function:")
				if len(s) == 2 {
					functionID = s[1]
				} else {
					panic("wrong format")
				}
			}

			createTaskParams := v1.CreateTaskParams{
				Name:       t.Name,
				Tags:       t.Tags,
				FunctionID: cache.Functions[functionID].ID,
				Interval:   t.Interval,
				ProjectID:  projectID,
			}
			taskObj, err := client.Client.Task.Create(&createTaskParams)
			if err != nil {
				panic(err)
			}
			color.Green("successfully created task %s", taskObj.ID)
			ct = taskObj
			cache.Tasks[ct.Name] = ct
		}
	}

	cache.Save()
	return nil
}

// buildSecrets will try to create or update all secrets in the current stack
func (b *Builder) buildSecrets() error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

	for _, s := range stack.Secrets {
		cg := cache.Secrets[s.Name]

		if cg != nil {
			color.Yellow("secret %s exists", s.Name)
		} else {
			createSecretParams := v1.CreateSecretParams{
				ProjectID: projectID,
				Value:     s.Value,
				Tags:      s.Tags,
				Name:      s.Name,
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
	cache.Save()
	return nil
}

// buildOAuth will try to create or update the oauth provider if it exists in the stack
func (b *Builder) buildOAuth() error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

	co := cache.OAuth

	if stack.OAuth.App.ClientName != "" {
		if co != nil {
			color.Yellow("oauth %s exists", co.ID)
			updateOAuthParams := v1.UpdateOAuthParams{
				OAuthID:      co.ID,
				RedirectUrls: stack.OAuth.App.RedirectUrls,
				LogoutUrls:   stack.OAuth.App.LogoutUrls,
			}
			oauthObj, err := client.Client.OAuth.Update(&updateOAuthParams)
			if err != nil {
				panic(err)
			}

			color.Green("successfully updated oauth %s", oauthObj.ID)

			co = oauthObj
			cache.OAuth = co
		} else {
			createOAuthParams := v1.CreateOAuthParams{
				ProjectID:    projectID,
				RedirectUrls: stack.OAuth.App.RedirectUrls,
				LogoutUrls:   stack.OAuth.App.LogoutUrls,
				AppName:      stack.OAuth.App.ClientName,
				Tags:         stack.OAuth.Tags,
			}

			oauthObj, err := client.Client.OAuth.Create(&createOAuthParams)
			if err != nil {
				panic(err)
			}

			color.Green("successfully created oauth %s", oauthObj.ID)

			co = oauthObj
			cache.OAuth = co
		}
	}
	cache.Save()
	return nil
}

// buildCollections will try to create or update all collections in the current stack
func (b *Builder) buildCollections() error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

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
	cache.Save()
	return nil
}

func (b *Builder) buildFiles() error {
	cache := b.cache
	stack := b.stack

	for _, filemapping := range stack.Filemapping {
		cacheFiles := cache.Filemappings[filemapping.Name]

		// check if we have old files
		fm := &Filemapping{
			Name:  filemapping.Name,
			Files: cacheFiles,
			Tags:  filemapping.Tags,
			Path:  filemapping.Path,
		}
		var files []*v1.File
		var err error

		// if cacheFiles is nil, we need to use the filemapping of the stack
		if cacheFiles == nil {
			files, err = filemapping.Upload(cache.Project.ID)
			if err != nil {
				log.Println(err)
				return err
			}
		} else {
			// otherwise use created filemapping with old files
			files, err = fm.Upload(cache.Project.ID)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		cache.Filemappings[filemapping.Name] = append(cache.Filemappings[filemapping.Name], files...)
	}
	cache.Save()
	return nil
}

func BuildStack() {
	builder := New()

	builder.buildSecrets()
	builder.buildOAuth()
	builder.buildCollections()
	builder.buildFunctions()
	builder.buildGateways()
	builder.buildTasks()
	builder.buildFiles()

	// verify that functions are running

}
