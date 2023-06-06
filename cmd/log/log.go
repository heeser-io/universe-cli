package log

import (
	"github.com/spf13/cobra"
)

var (
	LogCmd = &cobra.Command{
		Use:   "log",
		Short: "log api",
	}
)

func init() {
	LogCmd.AddCommand(ListCmd)
}
