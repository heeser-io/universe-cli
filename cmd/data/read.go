package data

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "reads a data object with the given id",
		Run: func(cmd *cobra.Command, args []string) {
			dataObj, err := client.Client.Data.Read(&v1.ReadDataParams{
				CollectionName: CollectionName,
				IndexName:      IndexName,
				IndexValue:     IndexValue,
				Filter:         Filter,
				AuthField:      AuthField,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(dataObj)))
		},
	}
)

func init() {

}
