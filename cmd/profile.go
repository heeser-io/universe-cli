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
	ProfileCmd.AddCommand(profile.UpdateCmd)
	ProfileCmd.AddCommand(profile.ListCmd)
	ProfileCmd.AddCommand(profile.CreateSubprofileCmd)
	ProfileCmd.AddCommand(profile.ReadSubprofileCmd)
	ProfileCmd.AddCommand(profile.AddGroupCmd)
	ProfileCmd.AddCommand(profile.RemoveGroupCmd)
}
