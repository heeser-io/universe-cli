package function

import (
	"github.com/spf13/cobra"
)

var (
	ProjectID   string
	FunctionID  string
	FunctionCmd = &cobra.Command{
		Use:   "function",
		Short: "function api",
	}
)

func init() {
	FunctionCmd.PersistentFlags().StringVarP(&ProjectID, "project-id", "p", "", "id of the project")
	FunctionCmd.PersistentFlags().StringVarP(&FunctionID, "function-id", "f", "", "id of the function")
	FunctionCmd.AddCommand(CreateCmd)
	FunctionCmd.AddCommand(ReadCmd)
	FunctionCmd.AddCommand(UpdateCmd)
	FunctionCmd.AddCommand(DeleteCmd)
	FunctionCmd.AddCommand(ListCmd)
}
