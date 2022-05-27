package builder

import (
	"github.com/heeser-io/universe-cli/config"
	v1 "github.com/heeser-io/universe/api/v1"
)

func HasChange() bool {
	cache := LoadOrCreate()

	cs := Checksum(config.Main.GetString("stackFile"))
	return cs != cache.Checksum
}
func Verify() {
	cache := LoadOrCreate()

	p, err := client.Project.Read(&v1.ReadProjectParams{
		ProjectID: cache.Project.ID,
	})
	if err != nil {
		panic(err)
	}
	cache.Project = p

	for _, secretObj := range cache.Secrets {
		readSecretParams := &v1.ReadSecretParams{
			SecretID: secretObj.ID,
		}

		s, err := client.Secret.Read(readSecretParams)
		if err != nil {
			panic(err)
		}
		cache.Secrets[s.Name] = s
	}

	for _, collectionObj := range cache.Collections {
		readCollectionParams := &v1.ReadCollectionParams{
			CollectionID: collectionObj.ID,
		}

		c, err := client.Collection.Read(readCollectionParams)
		if err != nil {
			panic(err)
		}
		cache.Collections[c.Name] = c
	}

	for _, functionObj := range cache.Functions {
		readFunctionParams := &v1.ReadFunctionParams{
			FunctionID: functionObj.ID,
		}

		f, err := client.Function.Read(readFunctionParams)
		if err != nil {
			panic(err)
		}

		cache.Functions[f.Name] = f
	}

	for _, gatewayObj := range cache.Gateways {
		readGatewayParams := &v1.ReadGatewayParams{
			GatewayID: gatewayObj.ID,
		}

		g, err := client.Gateway.Read(readGatewayParams)
		if err != nil {
			panic(err)
		}
		cache.Gateways[g.Name] = gatewayObj
	}
	cs := Checksum(config.Main.GetString("stackFile"))
	cache.Checksum = cs
	cache.Save()

}
