package task

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes a task with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.Task.Delete(&v1.DeleteTaskParams{
				TaskID: TaskID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			color.Green("successfully deleted task %s", TaskID)
		},
	}
)
