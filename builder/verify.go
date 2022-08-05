package builder

import (
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
)

func (b *Builder) HasChange() bool {
	cache := LoadOrCreate(b.path)

	cs, err := Checksum(GetStackFile(b.path))
	if err != nil {
		panic(err)
	}
	return cs != cache.Checksum
}
func (b *Builder) Verify() {
	cache := LoadOrCreate(b.path)

	p, err := client.Client.Project.Read(&v1.ReadProjectParams{
		ProjectID: cache.Project.ID,
	})
	if err != nil {
		panic(err)
	}
	cache.Project = p

	if cache.OAuth != nil {
		oauth, err := client.Client.OAuth.Read(&v1.ReadOAuthParams{
			OAuthID: cache.OAuth.ID,
		})
		if err != nil {
			panic(err)
		}
		cache.OAuth = oauth
	}

	for _, secretObj := range cache.Secrets {
		readSecretParams := &v1.ReadSecretParams{
			SecretID: secretObj.ID,
		}

		s, err := client.Client.Secret.Read(readSecretParams)
		if err != nil {
			panic(err)
		}
		cache.Secrets[s.Name] = s
	}

	for _, functionObj := range cache.Functions {
		readFunctionParams := &v1.ReadFunctionParams{
			FunctionID: functionObj.ID,
		}

		f, err := client.Client.Function.Read(readFunctionParams)
		if err != nil {
			panic(err)
		}

		cache.Functions[f.Name] = f
	}

	for _, gatewayObj := range cache.Gateways {
		readGatewayParams := &v1.ReadGatewayParams{
			GatewayID: gatewayObj.ID,
		}

		g, err := client.Client.Gateway.Read(readGatewayParams)
		if err != nil {
			panic(err)
		}
		cache.Gateways[g.Name] = gatewayObj
	}

	for _, taskObj := range cache.Tasks {
		readTaskParams := &v1.ReadTaskParams{
			TaskID: taskObj.ID,
		}

		t, err := client.Client.Task.Read(readTaskParams)
		if err != nil {
			panic(err)
		}
		cache.Tasks[t.Name] = t
	}

	cs, err := Checksum(GetStackFile(b.path))
	if err != nil {
		panic(err)
	}
	cache.Checksum = cs
	cache.Save()

}
