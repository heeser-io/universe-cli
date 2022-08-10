package file

import (
	"github.com/spf13/cobra"
)

var (
	Tags    *[]string
	FileID  string
	FileCmd = &cobra.Command{
		Use:   "file",
		Short: "file api",
	}
)

func init() {
	FileCmd.PersistentFlags().StringVarP(&FileID, "file-id", "p", "", "id of the file")
	FileCmd.AddCommand(ReadCmd)
	FileCmd.AddCommand(UpdateCmd)
	FileCmd.AddCommand(DeleteCmd)
	FileCmd.AddCommand(ListCmd)
}
