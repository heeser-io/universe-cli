package file

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "updates a file with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			fileObj, err := client.Client.File.Update(&v1.UpdateFileParams{
				FileID: FileID,
				Tags:   *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(fileObj)))
		},
	}
)

func init() {
	Tags = UpdateCmd.Flags().StringArray("tags", nil, "tags of the file")
}
