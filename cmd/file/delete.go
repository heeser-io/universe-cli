package file

import (
	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes a file with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Client.File.Delete(&v2.DeleteFileParams{
				FileID: FileID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			color.Green("successfully deleted file %s", FileID)
		},
	}
)
