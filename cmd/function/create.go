package function

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Name      string
	Tags      *[]string
	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "creates a function with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			functionObj, err := client.Client.Function.Create(&v1.CreateFunctionParams{
				ProjectID: ProjectID,
				Name:      Name,
				Tags:      *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(functionObj)))
		},
	}
)

func init() {
	CreateCmd.Flags().StringVar(&Name, "name", "", "name of the function")
	Tags = CreateCmd.Flags().StringArray("tags", nil, "tags of the function")
}
