package main

import (
	"errors"

	"github.com/heeser-io/universe-cli/client"
	"github.com/heeser-io/universe-cli/cmd"
	"github.com/heeser-io/universe-cli/config"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"github.com/thoas/go-funk"
)

func main() {
	godotenv.Load()
	config.Init()

	apiKey := config.Main.Get("apiKey")

	validate := func(input string) error {
		if len(input) != 64 {
			return errors.New("apikey length should be 64")
		}
		return nil
	}
	if funk.IsZero(apiKey) {
		prompt := promptui.Prompt{
			Label:    "apiKey",
			Validate: validate,
		}

		res, err := prompt.Run()
		if err != nil {
			panic(err)
		}
		config.Main.Set("apiKey", res)
		if err := config.Main.WriteConfig(); err != nil {
			panic(err)
		}
	}

	client.Init()
	cmd.Execute()
}
