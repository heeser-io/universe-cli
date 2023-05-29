package builder

import (
	"fmt"

	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
)

func (b *Builder) ListProjects() ([]v1.Project, error) {
	projects, err := client.Client.Project.List(&v1.ListProjectParams{})
	if err != nil {
		return nil, err
	}

	return projects, nil
}
func (b *Builder) Download(project *v1.Project) error {
	// we build our cache and save it so .stack

	filter := map[string]string{
		"projectId": project.ID,
	}

	c := Cache{
		Project:       project,
		Functions:     map[string]*v1.Function{},
		Secrets:       map[string]*v1.Secret{},
		Tasks:         map[string]*v1.Task{},
		Gateways:      map[string]*v1.Gateway{},
		Templates:     map[string]*v1.Template{},
		Subscriptions: map[string]*v1.Subscription{},
		Domains:       map[string]*v1.Domain{},
		Filemappings:  map[string][]*v1.File{},
		Webhooks:      map[string]*v1.Webhook{},
	}

	functions, err := client.Client.Function.List(&v1.ListFunctionParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, functionObj := range functions {
		c.Functions[functionObj.Name] = &functionObj
	}
	// c.Save()

	secrets, err := client.Client.Secret.List(&v1.ListSecretParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, secretObj := range secrets {
		c.Secrets[secretObj.Name] = &secretObj
	}

	gateways, err := client.Client.Gateway.List(&v1.ListGatewayParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, gatewayObj := range gateways {
		c.Gateways[gatewayObj.Name] = &gatewayObj
	}

	subscriptions, err := client.Client.Subscription.List(&v1.ListSubscriptionParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, subscriptionObj := range subscriptions {
		c.Subscriptions[subscriptionObj.Name] = &subscriptionObj
	}

	templates, err := client.Client.Template.List(&v1.ListTemplateParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, templateObj := range templates {
		c.Templates[templateObj.Name] = &templateObj
	}

	webhooks, err := client.Client.Webhook.List(&v1.ListWebhookParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, webhookObj := range webhooks {
		c.Webhooks[webhookObj.Name] = &webhookObj
	}

	tasks, err := client.Client.Task.List(&v1.ListTaskParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, taskObj := range tasks {
		c.Tasks[taskObj.Name] = &taskObj
	}

	domains, err := client.Client.Domain.List(&v1.ListDomainParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, domainObj := range domains {
		c.Domains[domainObj.Name] = &domainObj
	}

	oauth, err := client.Client.OAuth.List(&v1.ListOAuthParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	if len(oauth) > 0 {
		c.OAuth = &oauth[0]
	}

	keyvalues, err := client.Client.KeyValue.List(&v1.ListKeyValueParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, keyvalue := range keyvalues {
		c.KeyValues[fmt.Sprintf("%s.%s", keyvalue.Namespace, keyvalue.Key)] = &keyvalue
	}

	c.Save()
	return nil
}
