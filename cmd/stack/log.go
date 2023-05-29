package stack

import (
	"github.com/heeser-io/universe-cli/builder"
	"github.com/spf13/cobra"
)

var (
	LogCmd = &cobra.Command{
		Use:   "log",
		Short: "get live logs of stack",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := builder.New("", false)
			if err != nil {
				panic(err)
			}

			stack, err := cmd.Flags().GetString("stack")
			if err != nil {
				panic(err)
			}

			if stack == "" {
				b.Logs()
			}
		},
	}
)

func init() {
}
