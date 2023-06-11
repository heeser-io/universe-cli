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
	Name             string
	CreateItemString string
	CreateCmd        = &cobra.Command{
		Use:   "create",
		Short: "creates a data object with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			item := map[string]interface{}{}
			if err := json.Unmarshal([]byte(CreateItemString), &item); err != nil {
				panic(err)
			}

			dataObj, err := client.Client.Data.Create(&v1.CreateDataParams{
				CollectionName: CollectionName,
				Item:           item,
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
	CreateCmd.Flags().StringVar(&CreateItemString, "item-json", "", "item you want to create as json string")
}
