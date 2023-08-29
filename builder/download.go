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
		KeyValues:     map[string]*v1.KeyValue{},
	}

	functions, err := client.Client.Function.List(&v1.ListFunctionParams{
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

	secrets, err := client.Client.Secret.List(&v1.ListSecretParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, secretObj := range secrets {
		secretObjCopy := secretObj
		c.Secrets[secretObj.Name] = &secretObjCopy
	}

	gateways, err := client.Client.Gateway.List(&v1.ListGatewayParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, gatewayObj := range gateways {
		gatewayObjCopy := gatewayObj
		c.Gateways[gatewayObj.Name] = &gatewayObjCopy
	}

	subscriptions, err := client.Client.Subscription.List(&v1.ListSubscriptionParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, subscriptionObj := range subscriptions {
		subscriptionObjCopy := subscriptionObj
		c.Subscriptions[subscriptionObj.Name] = &subscriptionObjCopy
	}

	templates, err := client.Client.Template.List(&v1.ListTemplateParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, templateObj := range templates {
		templateObjCopy := templateObj
		c.Templates[templateObj.Name] = &templateObjCopy
	}

	webhooks, err := client.Client.Webhook.List(&v1.ListWebhookParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, webhookObj := range webhooks {
		webhookObjCopy := webhookObj
		c.Webhooks[webhookObj.Name] = &webhookObjCopy
	}

	tasks, err := client.Client.Task.List(&v1.ListTaskParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, taskObj := range tasks {
		taskObjCopy := taskObj
		c.Tasks[taskObj.Name] = &taskObjCopy
	}

	domains, err := client.Client.Domain.List(&v1.ListDomainParams{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	for _, domainObj := range domains {
		domainObjCopy := domainObj
		c.Domains[domainObj.Name] = &domainObjCopy
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
		keyvalueCopy := keyvalue
		c.KeyValues[fmt.Sprintf("%s.%s", keyvalue.Namespace, keyvalue.Key)] = &keyvalueCopy
	}

	c.Save()
	return nil
}
