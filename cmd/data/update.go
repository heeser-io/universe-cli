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
	UpdateItemString string
	UpdateCmd        = &cobra.Command{
		Use:   "update",
		Short: "updates a data object with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			updateValues := map[string]interface{}{}
			if err := json.Unmarshal([]byte(UpdateItemString), &updateValues); err != nil {
				panic(err)
			}

			dataObj, err := client.Client.Data.Update(&v1.UpdateDataParams{
				CollectionName: CollectionName,
				IndexName:      IndexName,
				IndexValue:     IndexValue,
				Values:         updateValues,
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
	UpdateCmd.Flags().StringVar(&UpdateItemString, "values-json", "", "values you want to update as json string")
}
