package cmd

import (
	"fmt"
	"os"

	"github.com/heeser-io/universe-cli/cmd/function"
	"github.com/heeser-io/universe-cli/cmd/project"
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(StackCmd)
	rootCmd.AddCommand(ProfileCmd)
	rootCmd.AddCommand(project.ProjectCmd)
	rootCmd.AddCommand(function.FunctionCmd)

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "whether to show debug logs or not")
}
