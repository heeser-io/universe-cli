package project

import (
	"github.com/spf13/cobra"
)

var (
	ProjectID  string
	ProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "project api",
	}
)

func init() {
	ProjectCmd.PersistentFlags().StringVarP(&ProjectID, "project-id", "p", "", "id of the project")
	ProjectCmd.AddCommand(CreateCmd)
	ProjectCmd.AddCommand(ReadCmd)
	ProjectCmd.AddCommand(UpdateCmd)
	ProjectCmd.AddCommand(DeleteCmd)
	ProjectCmd.AddCommand(ListCmd)
}
