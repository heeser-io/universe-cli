package websocket

import (
	"github.com/spf13/cobra"
)

var (
	WebsocketCmd = &cobra.Command{
		Use:   "websocket",
		Short: "websocket api",
	}
)

func init() {
	WebsocketCmd.AddCommand(ListCmd)
}
