package data

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	IndexField string
	Unique     bool
	IndexSort  string

	CreateIndexCmd = &cobra.Command{
		Use:   "create-index",
		Short: "creates an index on a collection",
		Run: func(cmd *cobra.Command, args []string) {
			dataObj, err := client.Client.Data.CreateIndex(&v2.CreateIndexParams{
				CollectionName: CollectionName,
				IndexField:     IndexField,
				Unique:         Unique,
				Sort:           IndexSort,
			})
			if err != nil {
				color.Red("err:%v\n", err)
				return
			}
			fmt.Println(string(v2.StructToByte(dataObj)))
		},
	}
)

func init() {
	CreateIndexCmd.Flags().StringVar(&IndexField, "index-field", "", "the field inside the collection to create the index on")
	CreateIndexCmd.Flags().BoolVar(&Unique, "unique", false, "the field inside the collection to create the index on")
	CreateIndexCmd.Flags().StringVar(&IndexSort, "sort", "asc", "define index sort")
}
