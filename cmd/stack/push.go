package stack

import (
	"github.com/heeser-io/universe-cli/builder"
	"github.com/spf13/cobra"
)

var (
	Force   bool
	PushCmd = &cobra.Command{
		Use:   "push",
		Short: "pushes the current stack",
		Run: func(cmd *cobra.Command, args []string) {
			v := builder.HasChange()

			if v || Force {
				builder.BuildStack()
				builder.Verify()
			}
		},
	}
)

func init() {
	PushCmd.Flags().BoolVarP(&Force, "force", "f", false, "force push")
}
