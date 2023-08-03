package file

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "reads a file with the given id",
		Run: func(cmd *cobra.Command, args []string) {
			fileObj, err := client.Client.File.Read(&v2.ReadFileParams{
				FileID: FileID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(fileObj)))
		},
	}
)
