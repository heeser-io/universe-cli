package file

import (
	"github.com/spf13/cobra"
)

var (
	Tags      *[]string
	FileID    string
	ProjectID string
	FileCmd   = &cobra.Command{
		Use:   "file",
		Short: "file api",
	}
)

func init() {
	FileCmd.PersistentFlags().StringVarP(&FileID, "file-id", "f", "", "id of the file")
	FileCmd.PersistentFlags().StringVarP(&FileID, "project-id", "p", "", "id of a project")
	FileCmd.AddCommand(ReadCmd)
	FileCmd.AddCommand(UpdateCmd)
	FileCmd.AddCommand(DeleteCmd)
	FileCmd.AddCommand(ListCmd)
	FileCmd.AddCommand(GenerateUrlCmd)
	FileCmd.AddCommand(UploadCmd)
}
