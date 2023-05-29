package stack

import (
	"os"

	"github.com/heeser-io/universe-cli/builder"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	StatusCmd = &cobra.Command{
		Use:   "status",
		Short: "list information about the stack",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := builder.New("./", false)
			if err != nil {
				panic(err)
			}

			stack, err := cmd.Flags().GetString("stack")
			if err != nil {
				panic(err)
			}

			t := table.NewWriter()

			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"Stack", "Resource", "Name", "Id"})

			if stack == "" {
				b.PrintStatus(&t)
			}

			subBuilders := b.GetSubBuilder()
			for _, subBuilder := range subBuilders {
				name := subBuilder.GetName()

				if name == stack || stack == "" {
					subBuilder.PrintStatus(&t)
				}
			}

			t.Render()
		},
	}
)
