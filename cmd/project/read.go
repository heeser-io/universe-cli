package project

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "reads a project with the given id",
		Run: func(cmd *cobra.Command, args []string) {
			projectObj, err := client.Client.Project.Read(&v1.ReadProjectParams{
				ProjectID: ProjectID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(projectObj)))
		},
	}
)
