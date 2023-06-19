package websocket

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	LeaveChannels    *[]string
	LeaveChannelsCmd = &cobra.Command{
		Use:   "leave-channels",
		Short: "leaves channels for a specific connectionId",
		Run: func(cmd *cobra.Command, args []string) {
			connections, err := client.Client.Websocket.LeaveChannels(&v1.LeaveChannelsParams{
				ConnectionID: ConnectionID,
				Channels:     *LeaveChannels,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(connections)))
		},
	}
)

func init() {
	LeaveChannels = LeaveChannelsCmd.Flags().StringSlice("channels", nil, "list of channels you want the connection to leave")
}
