package quota

import (
	"github.com/spf13/cobra"
)

var (
	TeamID   string
	QuotaCmd = &cobra.Command{
		Use:   "quota",
		Short: "quota api",
	}
)

func init() {
	QuotaCmd.AddCommand(ListCmd)
}
