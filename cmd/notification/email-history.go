package notification

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	Filter              map[string]string
	After               string
	AfterField          string
	Before              string
	BeforeField         string
	Sort                string
	SortOrder           string
	Limit               int64
	ListEmailHistoryCmd = &cobra.Command{
		Use:   "email-history",
		Short: "list email history for the current authenticated user",
		Run: func(cmd *cobra.Command, args []string) {
			emailHistory, err := client.Client.Email.ListHistory(&v2.ListEmailHistoryParams{
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
			fmt.Println(string(v2.StructToByte(emailHistory)))
		},
	}
)

func init() {
	ListEmailHistoryCmd.Flags().StringToStringVar(&Filter, "filter", map[string]string{}, "filters")
	ListEmailHistoryCmd.Flags().StringVar(&After, "after", "", "after")
	ListEmailHistoryCmd.Flags().StringVar(&AfterField, "afterField", "", "afterField")
	ListEmailHistoryCmd.Flags().StringVar(&Before, "before", "", "before")
	ListEmailHistoryCmd.Flags().StringVar(&BeforeField, "beforeField", "", "beforeField")
	ListEmailHistoryCmd.Flags().StringVar(&Sort, "sort", "", "sort")
	ListEmailHistoryCmd.Flags().StringVar(&SortOrder, "sort-order", "", "sort")
	ListEmailHistoryCmd.Flags().Int64Var(&Limit, "limit", 0, "limit")
}
