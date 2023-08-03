package function

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	Name      string
	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "updates a function with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			functionObj, err := client.Client.Function.Update(&v2.UpdateFunctionParams{
				FunctionID: FunctionID,
				Name:       Name,
				Tags:       *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(functionObj)))
		},
	}
)

func init() {
	UpdateCmd.Flags().StringVar(&Name, "name", "", "name of the function")
	Tags = UpdateCmd.Flags().StringArray("tags", nil, "tags of the function")
}
