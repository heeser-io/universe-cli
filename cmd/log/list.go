package log

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
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
			logs, err := client.Client.Log.List(&v2.ListLogParams{
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
			fmt.Println(string(v2.StructToByte(logs)))
		},
	}
)

func init() {
	ListCmd.Flags().StringToStringVar(&Filter, "filter", map[string]string{}, `you can filter for a specific function --filter "resource=function:xyz"`)
	ListCmd.Flags().StringVar(&After, "after", "", "define a value after which entries should be shown")
	ListCmd.Flags().StringVar(&AfterField, "afterField", "", "for example createdAt (default is _id)")
	ListCmd.Flags().StringVar(&Before, "before", "", "define a value before which entries should be shown")
	ListCmd.Flags().StringVar(&BeforeField, "beforeField", "", "for example createdAt (default is _id)")
	ListCmd.Flags().StringVar(&Sort, "sort", "", "for example createdAt")
	ListCmd.Flags().StringVar(&SortOrder, "sort-order", "", "DESC or ASC (default is asc)")
	ListCmd.Flags().Int64Var(&Limit, "limit", 0, "limit (server default is 50, max 500)")
}
