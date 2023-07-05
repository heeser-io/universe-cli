package builder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/heeser-io/universe-cli/client"
	"github.com/heeser-io/universe-cli/helper"
	"github.com/heeser-io/universe-cli/proxy"
	"github.com/heeser-io/universe-cli/shell"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/heeser-io/universe/pkg/serverless/invoker"
	"github.com/heeser-io/universe/services/gateway"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Builder struct {
	cache *Cache
	stack *Stack
	path  string
	ctx   context.Context
}

func (b *Builder) GetName() string {
	return path.Base(b.path)
}

// New creates a new Builder for the given path.
// If it is called on the root stack, call it with an empty string.
func New(path string, download bool) (*Builder, error) {
	logger := zerolog.New(os.Stdout).With().Int64("time", time.Now().Unix()).Str("path", path).Logger().Level(zerolog.ErrorLevel)
	loggerCtx := logger.WithContext(context.Background())

	stack, err := ReadStack(GetStackFile(path))
	if err != nil {
		return nil, err
	}

	logger.Debug().Msgf("successfully got stack from %s", path)

	cache := LoadOrCreate(path)

	if stack.Project.Branch != "" {
		client.Client.SetBranch(stack.Project.Branch)
	}

	if cache.Project == nil && !download {
		projectObj, err := client.Client.Project.Create(&v1.CreateProjectParams{
			Name: stack.Project.Name,
			Tags: stack.Project.Tags,
		})
		if err != nil {
			panic(err)
		}
		cache.Project = projectObj
		color.Green("successfully created project %s", projectObj.ID)
	} else if !download {
		params := &v1.UpdateProjectParams{
			ProjectID: cache.Project.ID,
			Name:      stack.Project.Name,
			Tags:      stack.Project.Tags,
		}
		if stack.Project.Settings != nil {
			params.Settings = *stack.Project.Settings
		}

		projectObj, err := client.Client.Project.Update(params)
		if err != nil {
			panic(err)
		}
		cache.Project = projectObj
		color.Green("successfully updated project %s", projectObj.ID)
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

	if cache.Secrets == nil {
		cache.Secrets = make(map[string]*v1.Secret)
	}

	if cache.Tasks == nil {
		cache.Tasks = make(map[string]*v1.Task)
	}

	if cache.Filemappings == nil {
		cache.Filemappings = make(map[string][]*v1.File)
	}

	if cache.Subscriptions == nil {
		cache.Subscriptions = make(map[string]*v1.Subscription)
	}

	cache.Save()

	return &Builder{
		cache: cache,
		stack: stack,
		path:  path,
		ctx:   loggerCtx,
	}, nil
}

// returns the projectID of the project
func (b *Builder) getProjectID() string {
	if b.cache == nil || b.cache.Project == nil {
		panic("no project in current stack")
	}
	return b.cache.Project.ID
}

func (b *Builder) GetStack() *Stack {
	return b.stack
}

func (b *Builder) Serve() error {
	stack := b.stack

	p := proxy.New()

	fPort := 19111

	fnPorts := map[string]int{}

	// Serve functions
	for _, function := range stack.Functions {
		// RUN FUNCTION
		if function.Language == "golang" {
			// RUN BINARY
			envStr := ""

			for key, value := range function.Environment {
				if strings.Contains(value, "secret:") {
					secretValue := funk.Find(stack.Secrets, func(x v1.Secret) bool {
						return x.Name == strings.Split(value, "secret:")[1]
					})

					secretObj, ok := secretValue.(v1.Secret)
					if !ok {
						panic(fmt.Errorf("secret value %s does not exist", value))
					}
					envStr += fmt.Sprintf("UNIVERSE_%s=%s ", key, secretObj.Value)
				} else {
					envStr += fmt.Sprintf("UNIVERSE_%s=%s ", key, value)
				}
			}
			helper.KillPort(fmt.Sprintf("%d", fPort))
			shell.CallAsync(fmt.Sprintf("%s RPC_PORT=%d go run %s/main.go", envStr, fPort, function.Name))
		}

		// Add to proxy
		p.AddFunction(function.Name, "localhost", fmt.Sprintf("%d", fPort))

		fnPorts[function.Name] = fPort
		// increment port
		fPort++
	}

	// now we serve all gateways
	for _, gateway := range stack.Gateways {
		router := mux.NewRouter()

		router.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				h.ServeHTTP(w, r)
			})
		})

		for _, route := range gateway.Routes {
			route := route
			fnPort := fnPorts[route.FunctionID]

			router.Methods(route.Method).Subrouter().HandleFunc(route.Path, CreateHandlerFunc(fnPort, route))

			p.Add(&gateway, route, fmt.Sprintf("localhost:%d", fPort))
		}

		go func(port int) {
			helper.KillPort(fmt.Sprintf("%d", fPort))
			if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), router); err != nil {
				panic(err)
			}
		}(fPort)
	}

	// serve all tasks
	c := cron.New()
	for _, task := range stack.Tasks {
		c.AddFunc(task.Interval, func() {
			fnName := strings.Split(task.FunctionID, "function:")[1]
			fnPort := fnPorts[fnName]

			conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", fnPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				panic(err)
			}

			rpcClient := invoker.NewInvokerClient(conn)

			rpcClient.Invoke(context.Background(), &invoker.InvokeParams{
				Auth:   &invoker.Auth{},
				Params: &invoker.RouteParams{},
			})
			if err != nil {
				panic(err)
			}
		})
	}

	c.Start()

	return p.Listen()
}

func CreateHandlerFunc(port int, route gateway.Route) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}

		rpcClient := invoker.NewInvokerClient(conn)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		query := map[string]string{}

		for k, v := range r.URL.Query() {
			query[k] = v[0]
		}

		headers := map[string]string{}

		for k, v := range r.Header {
			headers[k] = v[0]
		}

		profileObj, err := client.Client.Profile.Read(&v1.ReadProfileParams{})
		if err != nil {
			panic(err)
		}

		res, err := rpcClient.Invoke(context.Background(), &invoker.InvokeParams{
			Auth: &invoker.Auth{
				ApiKey:  client.ApiKey,
				Subject: profileObj.Subject,
			},
			Params: &invoker.RouteParams{
				Body:    string(b),
				Query:   query,
				Headers: headers,
				Params:  mux.Vars(r),
				Method:  r.Method,
			},
		})
		if err != nil {
			panic(err)
		}

		for k, v := range res.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(int(res.Status))
		if res.Body != "" {
			w.Write([]byte(res.Body))
		}
	}
}

// buildKeyValues will try to create or update all keyvalues in the current stack
func (b *Builder) buildKeyValues() error {
	cache := b.cache
	stack := b.stack
	projectID := b.getProjectID()

	for _, kv := range cache.KeyValues {
		// look for deleted ones
		deleteKeyValue := funk.Find(stack.KeyValues, func(keyvalue v1.KeyValue) bool {
			return keyvalue.Namespace == kv.Namespace && keyvalue.Key == kv.Key
		}) == nil

		if deleteKeyValue {
			deleteKeyvalueParams := v1.DeleteKeyValueParams{
				Namespace: kv.Namespace,
				Key:       kv.Key,
			}

			if err := client.Client.KeyValue.Delete(&deleteKeyvalueParams); err != nil {
				panic(err)
			}
			color.Red("keyvalue %s.%s (%s) deleted", kv.Namespace, kv.Key, kv.ID)
			delete(cache.KeyValues, fmt.Sprintf("%s.%s", kv.Namespace, kv.Key))
		}
	}

	for _, kv := range stack.KeyValues {
		key := fmt.Sprintf("%s.%s", kv.Namespace, kv.Key)
		ckv := cache.KeyValues[key]

		if ckv != nil {
			color.Yellow("keyvalue %s.%s exists", kv.Namespace, kv.Key)
			updateKeyvalueParams := v1.UpdateKeyValueParams{
				Key:       kv.Key,
				Namespace: kv.Namespace,
				Value:     kv.Value,
			}

			keyvalueObj, err := client.Client.KeyValue.Update(&updateKeyvalueParams)
			if err != nil {
				panic(err)
			}

			color.Green("successfully updated keyvalue %s.%s (%s)", kv.Namespace, kv.Key, keyvalueObj.ID)
			ckv = keyvalueObj
		} else {
			createKeyValue := v1.CreateKeyValueParams{
				ProjectID: projectID,
				Namespace: kv.Namespace,
				Key:       kv.Key,
				Value:     kv.Value,
				Tags:      kv.Tags,
			}

			keyvalueObj, err := client.Client.KeyValue.Create(&createKeyValue)
			if err != nil {
				panic(err)
			}

			color.Green("successfully created keyvalue %s.%s (%s)", kv.Namespace, kv.Key, keyvalueObj.ID)

			ckv = keyvalueObj
		}
		cache.KeyValues[key] = ckv
	}

	cache.Save()
	return nil
}

// buildFunctions will try to create or update all functions in the current stack
func (b *Builder) buildFunctions(functionName string) error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

	for _, s := range cache.Functions {
		// look for deleted ones
		deleteFunction := funk.Find(stack.Functions, func(function v1.Function) bool {
			return function.Name == s.Name
		}) == nil

		if deleteFunction {
			deleteFunctionParams := v1.DeleteFunctionParams{
				FunctionID: s.ID,
			}

			if err := client.Client.Function.Delete(&deleteFunctionParams); err != nil {
				panic(err)
			}
			color.Red("function %s (%s) deleted", s.Name, s.ID)
			delete(cache.Functions, s.Name)
		}
	}

	// push functions in parallel
	wg := sync.WaitGroup{}

	for _, function := range stack.Functions {
		if function.Name != functionName && functionName != "" {
			continue
		}

		wg.Add(1)
		go func(function v1.Function) error {
			lang := NewLanguage(function.Language)
			binaryBuilder := NewBinaryBuilder(lang, function, path.Join(stack.BasePath, b.path))
			_, err := binaryBuilder.Build()
			if err != nil {
				panic(err)
			}

			checksum, err := Checksum(path.Join(stack.BasePath, b.path, function.Path))
			if err != nil {
				return err
			}
			// functions
			cf := cache.Functions[function.Name]
			environment := function.Environment

			if funk.IsZero(environment) {
				environment = map[string]string{}
			}

			resources := []string{"secret", "template", "env"}

			for k, v := range environment {
				contains := false
				for _, resource := range resources {
					if strings.Contains(v, fmt.Sprintf("%s:", resource)) {
						s := strings.Split(v, ":")
						if len(s) == 2 {
							switch resource {
							case "secret":
								secretObj := cache.Secrets[s[1]]
								if secretObj != nil {
									environment[k] = fmt.Sprintf("secret:%s", secretObj.ID)
								}
							case "template":
								templateObj := cache.Templates[s[1]]
								if templateObj != nil {
									environment[k] = templateObj.ID
								}
							case "env":
								value := os.Getenv(s[1])
								if value == "" {
									panic(fmt.Sprintf("env var %s no value", s[1]))
								}
								environment[k] = value
							}
						}
						contains = true
					}
				}

				if !contains {
					environment[k] = v
				}
			}

			if cf != nil {
				// update file and function

				if checksum != cf.Checksum {
					functionObj, err := UpdateFunction(&v1.UpdateFunctionParams{
						FunctionID:  cf.ID,
						Path:        path.Join(stack.BasePath, b.path, function.Path),
						Checksum:    &checksum,
						Name:        function.Name,
						Environment: environment,
						Tags:        function.Tags,
						DockerArgs:  function.DockerArgs,
						Cache:       function.Cache,
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
					Path:        path.Join(stack.BasePath, b.path, function.Path),
					Handler:     function.Handler,
					Checksum:    checksum,
					ProjectID:   projectID,
					Permissions: function.Permissions,
					BaseImage:   function.BaseImage,
					Name:        function.Name,
					Tags:        function.Tags,
					DockerArgs:  function.DockerArgs,
					Language:    function.Language,
					Environment: environment,
					Cache:       function.Cache,
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
				if cache.Functions[r.FunctionID] == nil {
					panic(fmt.Sprintf("no function with name %s found", r.FunctionID))
				}
				routes = append(routes, gateway.Route{
					FunctionID:   cache.Functions[r.FunctionID].ID,
					Path:         r.Path,
					Method:       r.Method,
					AuthType:     r.AuthType,
					Groups:       r.Groups,
					ApiKeyHeader: r.ApiKeyHeader,
					Cache:        r.Cache,
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
				if cache.Functions[r.FunctionID] == nil {
					panic(fmt.Sprintf("no function with name %s found", r.FunctionID))
				}
				// r.FunctionID
				routes = append(routes, gateway.Route{
					FunctionID:   cache.Functions[r.FunctionID].ID,
					Path:         r.Path,
					Method:       r.Method,
					AuthType:     r.AuthType,
					Groups:       r.Groups,
					ApiKeyHeader: r.ApiKeyHeader,
					Cache:        r.Cache,
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
				Name:       t.Name,
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
				Reference:  t.Reference,
				Delay:      t.Delay,
				RunAt:      t.RunAt,
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

	for _, s := range cache.Secrets {
		// look for deleted ones
		deleteSecret := funk.Find(stack.Secrets, func(secret v1.Secret) bool {
			return secret.Name == s.Name
		}) == nil

		if deleteSecret {
			deleteSecretParams := v1.DeleteSecretParams{
				SecretID: s.ID,
			}

			if err := client.Client.Secret.Delete(&deleteSecretParams); err != nil {
				panic(err)
			}
			color.Red("secret %s (%s) deleted", s.Name, s.ID)
			delete(cache.Secrets, s.Name)
		}
	}

	for _, s := range stack.Secrets {
		cg := cache.Secrets[s.Name]

		secretValue := s.Value
		if strings.Contains(secretValue, "env:") {
			envSplit := strings.Split(secretValue, "env:")
			if len(envSplit) == 2 {
				envKey := envSplit[1]
				envVar := os.Getenv(envKey)

				if envVar == "" {
					return fmt.Errorf("env var %s is empty", envKey)
				}
				secretValue = envVar
			} else {
				return errors.New("should be in format env:VAR")
			}
		}

		if cg != nil {
			color.Yellow("secret %s exists", s.Name)
			updateSecretParams := v1.UpdateSecretParams{
				SecretID: cg.ID,
				Name:     s.Name,
				Tags:     s.Tags,
				Value:    secretValue,
			}

			secretObj, err := client.Client.Secret.Update(&updateSecretParams)
			if err != nil {
				panic(err)
			}

			color.Green("successfully updated secret %s", secretObj.ID)
			cg = secretObj
		} else {
			createSecretParams := v1.CreateSecretParams{
				ProjectID: projectID,
				Value:     secretValue,
				Tags:      s.Tags,
				Name:      s.Name,
			}

			secretObj, err := client.Client.Secret.Create(&createSecretParams)
			if err != nil {
				panic(err)
			}

			color.Green("successfully created secret %s", secretObj.ID)
			cg = secretObj
		}
		cache.Secrets[cg.Name] = cg
	}
	cache.Save()
	return nil
}

func (b *Builder) buildTemplates() error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

	for _, t := range stack.Templates {
		ct := cache.Templates[t.Name]

		body := t.Body

		splitBody := strings.Split(t.Body, "localfile:")

		if len(splitBody) == 2 {
			filepath := splitBody[1]

			wd, err := os.Getwd()
			if err != nil {
				log.Println(err)
				return err
			}
			b, err := os.ReadFile(path.Join(wd, filepath))
			if err != nil {
				log.Println(err)
				return err
			}

			body = string(b)
		}

		if ct != nil {
			color.Yellow("template %s exists", t.Name)
			updateTemplateParams := v1.UpdateTemplateParams{
				TemplateID: ct.ID,
				Name:       t.Name,
				Tags:       t.Tags,
				Body:       body,
			}

			templateObj, err := client.Client.Template.Update(&updateTemplateParams)
			if err != nil {
				return err
			}
			color.Green("successfully updated template %s", templateObj.Name)

			ct = templateObj
		} else {
			createTemplateParams := v1.CreateTemplateParams{
				ProjectID: projectID,
				Name:      t.Name,
				Tags:      t.Tags,
				Body:      body,
			}

			templateObj, err := client.Client.Template.Create(&createTemplateParams)
			if err != nil {
				return err
			}

			color.Green("successfully created template %s", templateObj.Name)

			ct = templateObj
		}
		cache.Templates[t.Name] = ct
		cache.Save()
	}

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
		}
		cache.OAuth = co
	}
	cache.Save()
	return nil
}

func (b *Builder) buildSubscriptions() error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

	for _, subscription := range stack.Subscriptions {
		cachedSubscription := cache.Subscriptions[subscription.Name]
		if cachedSubscription != nil {
			updateSubscriptionParams := v1.UpdateSubscriptionParams{
				SubscriptionID: cachedSubscription.ID,
				Name:           subscription.Name,
				Resource:       subscription.Resource,
				FunctionID:     cache.Functions[subscription.FunctionID].ID,
				Reference:      subscription.Reference,
			}

			subscriptionObj, err := client.Client.Subscription.Update(&updateSubscriptionParams)
			if err != nil {
				return err
			}
			color.Green("successfully updated subscription %s (%s)", subscriptionObj.Name, subscriptionObj.ID)
			cachedSubscription = subscriptionObj
		} else {
			createSubscriptionParams := v1.CreateSubscriptionParams{
				Name:       subscription.Name,
				ProjectID:  projectID,
				Resource:   subscription.Resource,
				FunctionID: cache.Functions[subscription.FunctionID].ID,
				Reference:  subscription.Reference,
			}

			subscriptionObj, err := client.Client.Subscription.Create(&createSubscriptionParams)
			if err != nil {
				return err
			}
			color.Green("successfully created subscription %s (%s)", subscriptionObj.Name, subscriptionObj.ID)
			cachedSubscription = subscriptionObj
		}
		cache.Subscriptions[subscription.Name] = cachedSubscription
	}

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
			Path:  path.Join(b.path, filemapping.Path),
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

func (b *Builder) buildDomains() error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

	for _, domain := range stack.Domains {
		cachedDomains := cache.Domains[domain.Name]

		gatewayID := domain.GatewayID

		if strings.Contains(gatewayID, "gateway:") {
			gatewayName := strings.Split(gatewayID, "gateway:")[1]
			gatewayID = cache.Gateways[gatewayName].ID
		}

		if cachedDomains != nil {
			updateDomainParams := v1.UpdateDomainParams{
				DomainID: cachedDomains.ID,
				Tags:     domain.Tags,
			}

			domainObj, err := client.Client.Domain.Update(&updateDomainParams)
			if err != nil {
				return err
			}
			color.Green("successfully updated domain %s (%s)", domainObj.Name, domainObj.ID)
			cachedDomains = domainObj
		} else {
			createDomainParams := v1.CreateDomainParams{
				Name:      domain.Name,
				GatewayID: gatewayID,
				Tags:      domain.Tags,
				ProjectID: projectID,
			}

			domainObj, err := client.Client.Domain.Create(&createDomainParams)
			if err != nil {
				return err
			}
			color.Green("successfully created domain %s (%s)", domainObj.Name, domainObj.ID)
			cachedDomains = domainObj
		}
		cache.Domains[domain.Name] = cachedDomains
	}
	cache.Save()
	return nil
}

func (b *Builder) buildWebhooks() error {
	cache := b.cache
	stack := b.stack

	projectID := b.getProjectID()

	for _, webhook := range stack.Webhooks {
		cachedWebhook := cache.Webhooks[webhook.Name]

		targetUrl := webhook.TargetUrl

		extractedVar := helper.ExtractVar(targetUrl)
		if extractedVar != nil {
			if extractedVar.Resource == "domain" {
				targetUrl = fmt.Sprintf("%s%s", cache.Domains[extractedVar.ID].Url, extractedVar.Appendix)
			}
		}
		if cachedWebhook != nil {
			updateWebhookParams := v1.UpdateWebhookParams{
				WebhookID: cachedWebhook.ID,
				Name:      webhook.Name,
				Resources: webhook.Resources,
				Receiver:  webhook.Receiver,
				TargetUrl: targetUrl,
				Tags:      webhook.Tags,
			}

			webhookObj, err := client.Client.Webhook.Update(&updateWebhookParams)
			if err != nil {
				return err
			}
			color.Green("successfully updated webhook %s (%s)", webhookObj.Name, webhookObj.ID)
			cachedWebhook = webhookObj
		} else {
			createWebhookParams := v1.CreateWebhookParams{
				Name:      webhook.Name,
				ProjectID: projectID,
				TargetUrl: targetUrl,
				Receiver:  webhook.Receiver,
				Resources: webhook.Resources,
			}

			webhookObj, err := client.Client.Webhook.Create(&createWebhookParams)
			if err != nil {
				return err
			}
			color.Green("successfully created webhook %s (%s)", webhookObj.Name, webhookObj.ID)
			cachedWebhook = webhookObj
		}
		cache.Webhooks[webhook.Name] = cachedWebhook
	}
	cache.Save()
	return nil
}

func (b *Builder) BuildStack(functionName string) error {
	if err := b.buildKeyValues(); err != nil {
		return err
	}

	if err := b.buildTemplates(); err != nil {
		return err
	}

	if err := b.buildSecrets(); err != nil {
		return err
	}

	if err := b.buildOAuth(); err != nil {
		return err
	}

	if err := b.buildFunctions(functionName); err != nil {
		return err
	}
	if err := b.buildGateways(); err != nil {
		return err
	}
	if err := b.buildSubscriptions(); err != nil {
		return err
	}
	if err := b.buildTasks(); err != nil {
		return err
	}
	if err := b.buildFiles(); err != nil {
		return err
	}

	if err := b.buildDomains(); err != nil {
		return err
	}

	if err := b.buildWebhooks(); err != nil {
		return err
	}

	return nil
}
