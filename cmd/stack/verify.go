package stack

import (
	"github.com/heeser-io/universe-cli/builder"
	"github.com/spf13/cobra"
)

var (
	VerifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "verifies the current stack",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := builder.New("", false)
			if err != nil {
				panic(err)
			}
			b.Verify()
		},
	}
)
