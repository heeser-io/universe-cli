package data

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	Filter      map[string]string
	ExistFilter map[string]string
	After       string
	AfterField  string
	Before      string
	BeforeField string
	Sort        string
	SortOrder   string
	Limit       int64
	ListCmd     = &cobra.Command{
		Use:   "list",
		Short: "list all data for the current authenticated user in the given collection",
		Run: func(cmd *cobra.Command, args []string) {
			// transform exists filter to map[string]bool
			transformedFilter := map[string]bool{}

			for k, v := range ExistFilter {
				bv, err := strconv.ParseBool(v)
				if err == nil {
					transformedFilter[k] = bv
				}
			}

			datas, err := client.Client.Data.List(&v2.ListDataParams{
				CollectionName: CollectionName,
				Filter:         Filter,
				Limit:          Limit,
				Sort:           Sort,
				SortOrder:      SortOrder == "asc",
				After:          After,
				AfterField:     AfterField,
				Before:         Before,
				BeforeField:    BeforeField,
				ExistFilter:    transformedFilter,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v2.StructToByte(datas)))
		},
	}
)

func init() {
	ListCmd.Flags().StringVar(&After, "after", "", "after")
	ListCmd.Flags().StringVar(&AfterField, "afterField", "", "afterField")
	ListCmd.Flags().StringVar(&Before, "before", "", "before")
	ListCmd.Flags().StringVar(&BeforeField, "beforeField", "", "beforeField")
	ListCmd.Flags().StringVar(&Sort, "sort", "", "sort")
	ListCmd.Flags().StringVar(&SortOrder, "sort-order", "", "sort")
	ListCmd.Flags().Int64Var(&Limit, "limit", 0, "limit")
	ListCmd.Flags().StringToStringVar(&ExistFilter, "exist-filter", map[string]string{}, "filters for existance (or neg) of fields")
}
