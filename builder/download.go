package builder

import (
	"fmt"

	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
)

func (b *Builder) ListProjects() ([]v2.Project, error) {
	projects, err := client.Client.Project.List(&v2.ListProjectParams{})
	if err != nil {
		return nil, err
	}

	return projects, nil
}
func (b *Builder) Download(project *v2.Project) error {
	// we build our cache and save it so .stack

	filter := map[string]string{
		"projectId": project.ID,
	}

	c := Cache{
		Project:       project,
		Functions:     map[string]*v2.Function{},
		Secrets:       map[string]*v2.Secret{},
		Tasks:         map[string]*v2.Task{},
		Gateways:      map[string]*v2.Gateway{},
		Templates:     map[string]*v2.Template{},
		Subscriptions: map[string]*v2.Subscription{},
		Domains:       map[string]*v2.Domain{},
		Filemappings:  map[string][]*v2.File{},
		Webhooks:      map[string]*v2.Webhook{},
	}

	functions, err := client.Client.Function.List(&v2.ListFunctionParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, functionObj := range functions {
		functionObjCopy := functionObj
		c.Functions[functionObj.Name] = &functionObjCopy
	}
	// c.Save()

	secrets, err := client.Client.Secret.List(&v2.ListSecretParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, secretObj := range secrets {
		secretObjCopy := secretObj
		c.Secrets[secretObj.Name] = &secretObjCopy
	}

	gateways, err := client.Client.Gateway.List(&v2.ListGatewayParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, gatewayObj := range gateways {
		gatewayObjCopy := gatewayObj
		c.Gateways[gatewayObj.Name] = &gatewayObjCopy
	}

	subscriptions, err := client.Client.Subscription.List(&v2.ListSubscriptionParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, subscriptionObj := range subscriptions {
		subscriptionObjCopy := subscriptionObj
		c.Subscriptions[subscriptionObj.Name] = &subscriptionObjCopy
	}

	templates, err := client.Client.Template.List(&v2.ListTemplateParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, templateObj := range templates {
		templateObjCopy := templateObj
		c.Templates[templateObj.Name] = &templateObjCopy
	}

	webhooks, err := client.Client.Webhook.List(&v2.ListWebhookParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, webhookObj := range webhooks {
		webhookObjCopy := webhookObj
		c.Webhooks[webhookObj.Name] = &webhookObjCopy
	}

	tasks, err := client.Client.Task.List(&v2.ListTaskParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, taskObj := range tasks {
		taskObjCopy := taskObj
		c.Tasks[taskObj.Name] = &taskObjCopy
	}

	domains, err := client.Client.Domain.List(&v2.ListDomainParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, domainObj := range domains {
		domainObjCopy := domainObj
		c.Domains[domainObj.Name] = &domainObjCopy
	}

	oauth, err := client.Client.OAuth.List(&v2.ListOAuthParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	if len(oauth) > 0 {
		c.OAuth = &oauth[0]
	}

	keyvalues, err := client.Client.KeyValue.List(&v2.ListKeyValueParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, keyvalue := range keyvalues {
		keyvalueCopy := keyvalue
		c.KeyValues[fmt.Sprintf("%s.%s", keyvalue.Namespace, keyvalue.Key)] = &keyvalueCopy
	}

	c.Save()
	return nil
}
