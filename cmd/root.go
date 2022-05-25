package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/heeser-io/universe-cli/config"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:           "universe",
		Short:         "CLI for universe",
		SilenceErrors: true,
	}
)

func Execute() {
	// load cluster specific conf
	// home, err := homedir.Dir()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// configPath := (path.Join(home, "/.meta-next"))
	apiKey := config.Main.GetString("apiKey")

	if apiKey == "" {
		log.Fatal("No api key")
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(StackCmd)
}
