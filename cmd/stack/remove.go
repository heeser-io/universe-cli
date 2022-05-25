package stack

import (
	"github.com/heeser-io/universe-cli/builder"
	"github.com/spf13/cobra"
)

var (
	RemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "removes the current stack",
		Run: func(cmd *cobra.Command, args []string) {
			builder.RemoveStack()
		},
	}
)
