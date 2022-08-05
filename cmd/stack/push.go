package stack

import (
	"fmt"

	"github.com/heeser-io/universe-cli/builder"
	"github.com/spf13/cobra"
)

var (
	Force   bool
	PushCmd = &cobra.Command{
		Use:   "push",
		Short: "pushes the current stack",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := builder.New("")
			if err != nil {
				panic(err)
			}
			v := b.HasChange()

			stack, err := cmd.Flags().GetString("stack")
			if err != nil {
				panic(err)
			}

			if v || Force {
				if stack == "" {
					if err := b.BuildStack(); err != nil {
						fmt.Println(err)
					}
					b.Verify()
				}

				subBuilders := b.GetSubBuilder()
				for _, subBuilder := range subBuilders {
					name := subBuilder.GetName()
					if name == stack || stack == "" {
						if err := subBuilder.BuildStack(); err != nil {
							fmt.Println(err)
						}
						subBuilder.Verify()
					}
				}
			}
		},
	}
)

func init() {
	PushCmd.Flags().BoolVarP(&Force, "force", "f", false, "force push")
}
