package function

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	Tags      *[]string
	CreateCmd = &cobra.Command{
		Use:   "create name [--tags]",
		Short: "creates a function with the given params",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			functionObj, err := client.Client.Function.Create(&v2.CreateFunctionParams{
				ProjectID: ProjectID,
				Name:      args[0],
				Tags:      *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(functionObj)))
		},
	}
)

func init() {
	Tags = CreateCmd.Flags().StringArray("tags", nil, "tags of the function")
}
