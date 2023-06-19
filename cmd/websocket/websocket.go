package websocket

import (
	"github.com/spf13/cobra"
)

var (
	ConnectionID string
	WebsocketCmd = &cobra.Command{
		Use:   "websocket",
		Short: "websocket api",
	}
)

func init() {
	WebsocketCmd.AddCommand(ListCmd)
	WebsocketCmd.AddCommand(JoinChannelsCmd)
	WebsocketCmd.AddCommand(LeaveChannelsCmd)

	WebsocketCmd.PersistentFlags().StringVar(&ConnectionID, "connection-id", "", "connection id")
}
