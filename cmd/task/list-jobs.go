package task

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	ListJobsCmd = &cobra.Command{
		Use:   "list-jobs",
		Short: "list all task jobs for the given task id",
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := client.Client.Task.ListJobs(&v2.ListJobParams{
				TaskID:      TaskID,
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
			fmt.Println(string(v2.StructToByte(tasks)))
		},
	}
)

func init() {
	ListJobsCmd.Flags().StringToStringVar(&Filter, "filter", map[string]string{}, "filters")
	ListJobsCmd.Flags().StringVar(&After, "after", "", "after")
	ListJobsCmd.Flags().StringVar(&AfterField, "afterField", "", "afterField")
	ListJobsCmd.Flags().StringVar(&Before, "before", "", "before")
	ListJobsCmd.Flags().StringVar(&BeforeField, "beforeField", "", "beforeField")
	ListJobsCmd.Flags().StringVar(&Sort, "sort", "", "sort")
	ListJobsCmd.Flags().StringVar(&SortOrder, "sort-order", "", "sort")
	ListJobsCmd.Flags().Int64Var(&Limit, "limit", 0, "limit")
}
