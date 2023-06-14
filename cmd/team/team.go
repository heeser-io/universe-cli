package team

import (
	"github.com/spf13/cobra"
)

var (
	TeamCmd = &cobra.Command{
		Use:   "team",
		Short: "team api",
	}
)

func init() {
	TeamCmd.AddCommand(ReadAuthCmd)
}
