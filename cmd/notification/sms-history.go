package notification

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	ListSMSHistoryCmd = &cobra.Command{
		Use:   "sms-history",
		Short: "list sms history for the current authenticated user",
		Run: func(cmd *cobra.Command, args []string) {
			smsHistory, err := client.Client.SMS.ListHistory(&v1.ListSMSHistoryParams{
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
			fmt.Println(string(v1.StructToByte(smsHistory)))
		},
	}
)

func init() {
	ListSMSHistoryCmd.Flags().StringToStringVar(&Filter, "filter", map[string]string{}, "filters")
	ListSMSHistoryCmd.Flags().StringVar(&After, "after", "", "after")
	ListSMSHistoryCmd.Flags().StringVar(&AfterField, "afterField", "", "afterField")
	ListSMSHistoryCmd.Flags().StringVar(&Before, "before", "", "before")
	ListSMSHistoryCmd.Flags().StringVar(&BeforeField, "beforeField", "", "beforeField")
	ListSMSHistoryCmd.Flags().StringVar(&Sort, "sort", "", "sort")
	ListSMSHistoryCmd.Flags().StringVar(&SortOrder, "sort-order", "", "sort")
	ListSMSHistoryCmd.Flags().Int64Var(&Limit, "limit", 0, "limit")
}
