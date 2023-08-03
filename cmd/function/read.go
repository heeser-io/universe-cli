package function

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "reads a function with the given id",
		Run: func(cmd *cobra.Command, args []string) {
			functionObj, err := client.Client.Function.Read(&v2.ReadFunctionParams{
				FunctionID: FunctionID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(functionObj)))
		},
	}
)
