package file

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	Filter    map[string]string
	After     string
	Sort      string
	SortOrder string
	Limit     int64
	ListCmd   = &cobra.Command{
		Use:   "list",
		Short: "list all files for the current authenticated user",
		Run: func(cmd *cobra.Command, args []string) {
			files, err := client.Client.File.List(&v1.ListFileParams{
				Filter:    Filter,
				Limit:     Limit,
				Sort:      Sort,
				SortOrder: SortOrder == "asc",
				After:     After,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(files)))
		},
	}
)

func init() {
	ListCmd.Flags().StringToStringVar(&Filter, "filter", map[string]string{}, "filters")
	ListCmd.Flags().StringVar(&After, "after", "", "after")
	ListCmd.Flags().StringVar(&Sort, "sort", "", "sort")
	ListCmd.Flags().StringVar(&SortOrder, "sort-order", "", "sort")
	ListCmd.Flags().Int64Var(&Limit, "limit", 0, "limit")
}
