package data

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	DeleteIndexName string
	DeleteIndexCmd  = &cobra.Command{
		Use:   "delete-index",
		Short: "deletes an index from a collection",
		Run: func(cmd *cobra.Command, args []string) {
			dataObj, err := client.Client.Data.DeleteIndex(&v1.DeleteIndexParams{
				CollectionName: CollectionName,
				IndexName:      IndexName,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			fmt.Println(string(v1.StructToByte(dataObj)))
		},
	}
)

func init() {
}
