package quota

import (
	"fmt"

	"github.com/heeser-io/universe-cli/client"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/spf13/cobra"
)

var (
	ListCmd = &cobra.Command{
		Use:   "get-overview",
		Short: "returns an overview of quotas",
		Run: func(cmd *cobra.Command, args []string) {
			quotaOverview, err := client.Client.Quota.GetOverview(&v2.GetQuotaOverviewParams{})
			if err != nil {
				panic(err)
			}

			fmt.Println(string(v2.StructToByte(quotaOverview)))
		},
	}
)

func init() {
}
