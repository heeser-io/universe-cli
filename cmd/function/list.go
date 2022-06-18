package function

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "list all functions for the current authenticated user",
		Run: func(cmd *cobra.Command, args []string) {
			functions, err := client.Client.Function.List(&v1.ListFunctionParams{})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(functions)))
		},
	}
)
