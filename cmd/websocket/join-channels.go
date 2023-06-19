package websocket

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/heeser-io/universe-cli/client"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/spf13/cobra"
)

var (
	JoinChannels    *[]string
	JoinChannelsCmd = &cobra.Command{
		Use:   "join-channels",
		Short: "joins channels for a specific connectionId",
		Run: func(cmd *cobra.Command, args []string) {
			connections, err := client.Client.Websocket.JoinChannels(&v1.JoinChannelsParams{
				ConnectionID: ConnectionID,
				Channels:     *JoinChannels,
			})
			if err != nil {
				color.Red("err:%v\n", err)
			}
			fmt.Println(string(v1.StructToByte(connections)))
		},
	}
)

func init() {
	JoinChannels = JoinChannelsCmd.Flags().StringSlice("channels", nil, "list of channels you want the connection to join")
}
