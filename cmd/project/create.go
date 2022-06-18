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
		Run: func(cmd *cobra.Command, args []string) {
			projectObj, err := client.Client.Project.Create(&v1.CreateProjectParams{
				Name: Name,
				Tags: *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(projectObj)))
		},
	}
)

func init() {
	CreateCmd.Flags().StringVar(&Name, "name", "", "name of the project")
	Tags = CreateCmd.Flags().StringArray("tags", nil, "tags of the project")
}
