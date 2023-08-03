package data

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	PushArrayCmd = &cobra.Command{
		Use:   "push-array-element",
		Short: "pushes a element to a given array of a data object",
		Run: func(cmd *cobra.Command, args []string) {
			pushValue := map[string]interface{}{}
			if err := json.Unmarshal([]byte(UpdateItemString), &pushValue); err != nil {
				panic(err)
			}

			dataObj, err := client.Client.Data.Push(&v2.PushDataParams{
				CollectionName: CollectionName,
				IndexName:      IndexName,
				IndexValue:     IndexValue,
				Value:          pushValue,
				Key:            Key,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(dataObj)))
		},
	}
)

func init() {
	PushArrayCmd.Flags().StringVar(&Key, "key", "", "key of the array (e.g elements or groups)")
	PushArrayCmd.Flags().StringVar(&UpdateItemString, "value-json", "", "value you want to push as json string")
}
