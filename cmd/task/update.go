package task

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
		Short: "updates a task with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			taskObj, err := client.Client.Task.Update(&v2.UpdateTaskParams{
				TaskID: TaskID,
				Name:   Name,
				Tags:   *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(taskObj)))
		},
	}
)

func init() {
	UpdateCmd.Flags().StringVar(&Name, "name", "", "name of the task")
	Tags = UpdateCmd.Flags().StringArray("tags", nil, "tags of the task")
}
