package task

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Name       string
	Tags       *[]string
	FunctionID string
	Interval   string
	RunAt      string
	Delay      int64
	CreateCmd  = &cobra.Command{
		Use:   "create",
		Short: "creates a task with the given params",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskObj, err := client.Client.Task.Create(&v1.CreateTaskParams{
				FunctionID: FunctionID,
				Interval:   Interval,
				RunAt:      RunAt,
				Delay:      Delay,
				Name:       args[0],
				Tags:       *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			fmt.Println(string(v1.StructToByte(taskObj)))
		},
	}
)

func init() {
	Tags = CreateCmd.Flags().StringArray("tags", nil, "tags of the task")
}
