package log

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Filter      map[string]string
	After       string
	AfterField  string
	Before      string
	BeforeField string
	Sort        string
	SortOrder   string
	Limit       int64
	ListCmd     = &cobra.Command{
		Use:   "list",
		Short: "list all logs for the current authenticated user",
		Run: func(cmd *cobra.Command, args []string) {
			logs, err := client.Client.Log.List(&v1.ListLogParams{
				Filter:      Filter,
				Limit:       Limit,
				Sort:        Sort,
				SortOrder:   SortOrder == "asc",
				After:       After,
				AfterField:  AfterField,
				Before:      Before,
				BeforeField: BeforeField,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(logs)))
		},
	}
)

func init() {
	ListCmd.Flags().StringToStringVar(&Filter, "filter", map[string]string{}, "filters")
	ListCmd.Flags().StringVar(&After, "after", "", "after")
	ListCmd.Flags().StringVar(&AfterField, "afterField", "", "afterField")
	ListCmd.Flags().StringVar(&Before, "before", "", "before")
	ListCmd.Flags().StringVar(&BeforeField, "beforeField", "", "beforeField")
	ListCmd.Flags().StringVar(&Sort, "sort", "", "sort")
	ListCmd.Flags().StringVar(&SortOrder, "sort-order", "", "sort")
	ListCmd.Flags().Int64Var(&Limit, "limit", 0, "limit")
}
