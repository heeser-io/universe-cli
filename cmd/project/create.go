package project

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
		Short: "creates a project with the given params",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectObj, err := client.Client.Project.Create(&v1.CreateProjectParams{
				Name: args[0],
				Tags: *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			fmt.Println(string(v1.StructToByte(projectObj)))
		},
	}
)

func init() {
	Tags = CreateCmd.Flags().StringArray("tags", nil, "tags of the project")
}
