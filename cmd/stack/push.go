package stack

import (
	"github.com/heeser-io/universe-cli/builder"
	"github.com/spf13/cobra"
)

var (
	PushCmd = &cobra.Command{
		Use:   "push",
		Short: "pushes the current stack",
		Run: func(cmd *cobra.Command, args []string) {
			v := builder.HasChange()

			if v {
				builder.BuildStack()
				builder.Verify()
			}
		},
	}
)
