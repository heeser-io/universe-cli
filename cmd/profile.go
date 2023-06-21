package cmd

import (
	"github.com/heeser-io/universe-cli/cmd/profile"
	"github.com/spf13/cobra"
)

var (
	ProfileCmd = &cobra.Command{
		Use:   "profile",
		Short: "profile api",
	}
)

func init() {
	ProfileCmd.AddCommand(profile.ReadCmd)
	ProfileCmd.AddCommand(profile.ListCmd)
}
