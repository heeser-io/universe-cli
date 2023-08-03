package file

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	GenerateUrlCmd = &cobra.Command{
		Use:   "generate-url",
		Short: "generates a url for the specified file",
		Run: func(cmd *cobra.Command, args []string) {
			urlRes, err := client.Client.File.GenerateUrl(&v2.GenerateUrlParams{
				FileID: FileID,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(urlRes)))
		},
	}
)
