package data

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	Key              string
	PathKeys         *[]string
	PathValues       []interface{}
	PathValuesString *[]string
	UpdateArrayCmd   = &cobra.Command{
		Use:   "update-array-element",
		Short: "updates a data objects array with the given params",
		Run: func(cmd *cobra.Command, args []string) {
			updateValues := map[string]interface{}{}
			if err := json.Unmarshal([]byte(UpdateItemString), &updateValues); err != nil {
				panic(err)
			}

			// Transform pathvalues object
			for _, pathValueStr := range *PathValuesString {
				if i, err := strconv.Atoi(pathValueStr); err == nil {
					PathValues = append(PathValues, i)
					continue
				}

				if b, err := strconv.ParseBool(pathValueStr); err == nil {
					PathValues = append(PathValues, b)
					continue
				}

				if f, err := strconv.ParseFloat(pathValueStr, 64); err == nil {
					PathValues = append(PathValues, f)
					continue
				}

				PathValues = append(PathValues, pathValueStr)
			}

			dataObj, err := client.Client.Data.UpdateArray(&v2.UpdateDataArrayParams{
				CollectionName: CollectionName,
				IndexName:      IndexName,
				IndexValue:     IndexValue,
				Values:         updateValues,
				Key:            Key,
				PathKeys:       *PathKeys,
				PathValues:     PathValues,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(dataObj)))
		},
	}
)

func init() {
	UpdateArrayCmd.Flags().StringVar(&Key, "key", "", "key of the array (e.g elements or groups)")
	PathKeys = UpdateArrayCmd.Flags().StringSlice("path-keys", nil, "identifiers of the element inside the array you want to update (e.g id)")
	PathValuesString = UpdateArrayCmd.Flags().StringSlice("path-values", nil, "values you you want to update on the found element")
	UpdateArrayCmd.Flags().StringVar(&UpdateItemString, "values-json", "", "values you want to update as json string")
}
