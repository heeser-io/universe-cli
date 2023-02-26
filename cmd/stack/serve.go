package stack

import (
	"github.com/heeser-io/universe-cli/builder"
	"github.com/spf13/cobra"
)

var (
	ServeCmd = &cobra.Command{
		Use:   "serve",
		Short: "serves all functions for local testing",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := builder.New("")
			if err != nil {
				panic(err)
			}

			// stack, err := cmd.Flags().GetString("stack")
			// if err != nil {
			// 	panic(err)
			// }

			b.Serve()
		},
	}
)

func init() {
}
