package task

import (
	"github.com/spf13/cobra"
)

var (
	TaskID  string
	TaskCmd = &cobra.Command{
		Use:   "task",
		Short: "task api",
	}
)

func init() {
	TaskCmd.PersistentFlags().StringVarP(&TaskID, "task-id", "p", "", "id of the task")
	TaskCmd.AddCommand(CreateCmd)
	TaskCmd.AddCommand(ReadCmd)
	TaskCmd.AddCommand(UpdateCmd)
	TaskCmd.AddCommand(DeleteCmd)
	TaskCmd.AddCommand(ListCmd)
	TaskCmd.AddCommand(ListJobsCmd)
}
