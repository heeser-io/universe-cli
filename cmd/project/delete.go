package project

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes a project with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Project.Delete(&v1.DeleteProjectParams{
				ProjectID: ProjectID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			color.Green("successfully deleted project %s", ProjectID)
		},
	}
)
