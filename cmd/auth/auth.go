package auth

import "github.com/spf13/cobra"

var (
	AuthCmd = &cobra.Command{
		Use:   "auth",
		Short: "auth api",
	}
)

func init() {
	AuthCmd.AddCommand(MeCmd)
}
