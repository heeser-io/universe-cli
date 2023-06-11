package data

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	FindOrCreateCmd = &cobra.Command{
		Use:   "find-or-create",
		Short: "find or creates a data object with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			item := map[string]interface{}{}
			if err := json.Unmarshal([]byte(CreateItemString), &item); err != nil {
				panic(err)
			}

			dataObj, err := client.Client.Data.FindOrCreate(&v1.FindOrCreateParams{
				CollectionName: CollectionName,
				IndexName:      IndexName,
				IndexValue:     IndexValue,
				Item:           item,
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
	FindOrCreateCmd.Flags().StringVar(&CreateItemString, "item-json", "", "item you want to create as json string in case no item found")
}
