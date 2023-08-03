package file

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
		Short: "updates a file with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			fileObj, err := client.Client.File.Update(&v2.UpdateFileParams{
				FileID: FileID,
				Tags:   *Tags,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(fileObj)))
		},
	}
)

func init() {
	Tags = UpdateCmd.Flags().StringArray("tags", nil, "tags of the file")
}
