package client

import (
	"github.com/heeser-io/universe-cli/config"
	v2 "github.com/heeser-io/universe/api/v2"
)

var (
	Client *v2.Client
	ApiKey string
)

func Init() {
	ApiKey = config.Main.GetString("apiKey")

	token := config.Main.GetString("token")

	if token != "" {
		Client = v2.WithToken(token)
	} else {
		Client = v2.WithAPIKey(ApiKey)
	}
}
