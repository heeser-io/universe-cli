package job

import (
	"github.com/spf13/cobra"
)

var (
	JobID  string
	JobCmd = &cobra.Command{
		Use:   "job",
		Short: "job api",
	}
)

func init() {
	JobCmd.PersistentFlags().StringVarP(&JobID, "job-id", "p", "", "id of the job")
	JobCmd.AddCommand(ReadCmd)
	JobCmd.AddCommand(ListCmd)
}
