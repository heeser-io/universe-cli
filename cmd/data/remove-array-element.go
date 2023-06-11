package data

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Index                  int64
	RemovePathKeys         *[]string
	RemovePathValuesString *[]string
	RemoveArrayElementCmd  = &cobra.Command{
		Use:   "remove-array-element",
		Short: "removes an element from a given array inside a data object",
		Run: func(cmd *cobra.Command, args []string) {
			params := &v1.RemoveDataArrayElementParams{
				CollectionName: CollectionName,
				IndexName:      IndexName,
				IndexValue:     IndexValue,
				Key:            Key,
			}

			if Index != -1 {
				params.Index = &Index
			} else {
				params.PathKeys = *RemovePathKeys

				// Transform pathvalues object
				for _, pathValueStr := range *RemovePathValuesString {
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
				params.PathValues = PathValues
			}

			dataObj, err := client.Client.Data.RemoveArrayElement(params)
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(dataObj)))
		},
	}
)

func init() {
	RemoveArrayElementCmd.Flags().StringVar(&Key, "key", "", "key of the array (e.g elements or groups)")
	RemoveArrayElementCmd.Flags().Int64Var(&Index, "index", -1, "index of the element inside the array")
	RemovePathKeys = RemoveArrayElementCmd.Flags().StringSlice("path-keys", nil, "identifiers of the element inside the array you want to update (e.g id)")
	RemovePathValuesString = RemoveArrayElementCmd.Flags().StringSlice("path-values", nil, "values you you want to update on the found element")
}
