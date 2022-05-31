package client

import (
	"github.com/heeser-io/universe-cli/config"
	v1 "github.com/heeser-io/universe/api/v1"
)

var (
	Client *v1.Client
)

func Init() {
	apiKey := config.Main.GetString("apiKey")
	Client = v1.WithAPIKey(apiKey)
}
