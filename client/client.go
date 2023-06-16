package client

import (
	"github.com/heeser-io/universe-cli/config"
	v1 "github.com/heeser-io/universe/api/v1"
)

var (
	Client *v1.Client
	ApiKey string
)

func Init() {
	ApiKey = config.Main.GetString("apiKey")

	token := config.Main.GetString("token")

	if token != "" {
		Client = v1.WithToken(token)
	} else {
		Client = v1.WithAPIKey(ApiKey)
	}
}
