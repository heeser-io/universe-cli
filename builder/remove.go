package builder

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
)

func (b *Builder) RemoveStack() {
	cache := LoadOrCreate(b.path)

	if err := client.Client.Project.Delete(&v2.DeleteProjectParams{
		ProjectID: cache.Project.ID,
	}); err != nil {
		color.Red("no project with %s ", cache.Project.ID)
	} else {
		color.Green("successfully deleted project %s (%s)", cache.Project.Name, cache.Project.ID)
	}

	if cache.OAuth != nil {
		// remove oauth
		if err := client.Client.OAuth.Delete(&v2.DeleteOAuthParams{
			OAuthID: cache.OAuth.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted oauth %s (%s)", cache.OAuth.App.ClientName, cache.OAuth.ID)
	}

	// remove secrets
	for _, s := range cache.Secrets {
		if err := client.Client.Secret.Delete(&v2.DeleteSecretParams{
			SecretID: s.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted secret %s (%s)", s.ID, s.Name)
	}

	// remove functions
	for _, f := range cache.Functions {
		if err := client.Client.Function.Delete(&v2.DeleteFunctionParams{
			FunctionID: f.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted function %s (%s)", f.ID, f.Name)
	}

	// remove gateways
	for _, g := range cache.Gateways {
		if err := client.Client.Gateway.Delete(&v2.DeleteGatewayParams{
			GatewayID: g.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted gateway %s (%s)", g.ID, g.Name)
	}

	for _, t := range cache.Tasks {
		if err := client.Client.Task.Delete(&v2.DeleteTaskParams{
			TaskID: t.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted task %s (%s)", t.ID, t.Name)
	}

	for _, w := range cache.Webhooks {
		if err := client.Client.Webhook.Delete(&v2.DeleteWebhookParams{
			WebhookID: w.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted webhook %s (%s)", w.ID, w.Name)
	}

	for _, d := range cache.Domains {
		if err := client.Client.Domain.Delete(&v2.DeleteDomainParams{
			DomainID: d.ID,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted domain %s (%s)", d.ID, d.Name)
	}

	for _, kv := range cache.KeyValues {
		if err := client.Client.KeyValue.Delete(&v2.DeleteKeyValueParams{
			Namespace: kv.Namespace,
			Key:       kv.Key,
		}); err != nil {
			panic(err)
		}
		color.Green("successfully deleted keyvalue %s.%s (%s)", kv.Namespace, kv.Key, kv.ID)
	}

	// make a backup of .stack.yml for later recovery
}
