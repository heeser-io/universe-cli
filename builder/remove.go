package builder

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
)

func RemoveStack() {
	cache := LoadOrCreate()

	if err := client.Client.Project.Delete(&v1.DeleteProjectParams{
		ProjectID: cache.Project.ID,
	}); err != nil {
		panic(err)
	}

	color.Green("successfully deleted project %s", cache.Project.ID)

	if cache.OAuth != nil {
		// remove oauth
		if err := client.Client.OAuth.Delete(&v1.DeleteOAuthParams{
			OAuthID: cache.OAuth.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted oauth %s", cache.Project.ID)
	}

	// remove secrets
	for _, s := range cache.Secrets {
		if err := client.Client.Secret.Delete(&v1.DeleteSecretParams{
			SecretID: s.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted secret %s", s.ID)
	}

	// remove collections
	for _, c := range cache.Collections {
		if err := client.Client.Collection.Delete(&v1.DeleteCollectionParams{
			CollectionID: c.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted collection %s", c.ID)
	}

	// remove functions
	for _, f := range cache.Functions {
		if err := client.Client.Function.Delete(&v1.DeleteFunctionParams{
			FunctionID: f.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted function %s", f.ID)
	}

	// remove gateways
	for _, g := range cache.Gateways {
		if err := client.Client.Gateway.Delete(&v1.DeleteGatewayParams{
			GatewayID: g.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted gateway %s", g.ID)
	}

	// make a backup of .stack.yml for later recovery
}
