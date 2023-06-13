package notification

import (
	"github.com/spf13/cobra"
)

var (
	NotificationCmd = &cobra.Command{
		Use:   "notification",
		Short: "notification api",
	}
)

func init() {
	NotificationCmd.AddCommand(SendEmailCmd)
	NotificationCmd.AddCommand(SendSMSCmd)
	NotificationCmd.AddCommand(ListEmailHistoryCmd)
	NotificationCmd.AddCommand(ListSMSHistoryCmd)
}
