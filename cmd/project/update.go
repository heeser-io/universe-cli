package project

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "updates a project with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			projectObj, err := client.Client.Project.Update(&v2.UpdateProjectParams{
				ProjectID: ProjectID,
				Name:      Name,
				Tags:      *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(projectObj)))
		},
	}
)

func init() {
	UpdateCmd.Flags().StringVar(&Name, "name", "", "name of the project")
	Tags = UpdateCmd.Flags().StringArray("tags", nil, "tags of the project")
}
