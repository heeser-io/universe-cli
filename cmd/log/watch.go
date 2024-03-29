package log

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	ws "github.com/gorilla/websocket"
	"github.com/heeser-io/universe-cli/builder"
	"github.com/heeser-io/universe-cli/client"
	"github.com/heeser-io/universe-cli/config"
	v1 "github.com/heeser-io/universe/api/v1"
	"github.com/heeser-io/universe/services/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

var (
	WEBSOCKET_URL = "api.universecloud.io"
	Resources     *[]string
	WatchCmd      = &cobra.Command{
		Use:   "watch",
		Short: "watch live logs for the given resources",
		Run: func(cmd *cobra.Command, args []string) {
			if len(*Resources) == 0 {
				c := builder.LoadOrCreate("")
				resources := []string{}

				for _, function := range c.Functions {
					resources = append(resources, fmt.Sprintf("function:%s", function.ID))
				}
				(*Resources) = resources
			}

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt)

			u := url.URL{Scheme: "ws", Host: WEBSOCKET_URL, Path: "/websocket", RawQuery: fmt.Sprintf("devkey=%s", config.Main.GetString("apiKey"))}

			c, _, err := ws.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Fatal("dial:", err)
			}

			done := make(chan struct{})

			connection := make(chan websocket.Connection, 1)

			go func() {
				defer close(done)
				for {
					message := websocket.Message{}
					if err := c.ReadJSON(&message); err != nil {
						cerr, ok := err.(*ws.CloseError)
						if ok {
							log.Println("close error", cerr)
							return
						}
					}
					// if strings.Contains(string(message), "getConnection") {
					// 	messageObj
					// 	if err := json.Unmarshal(message, &connectionObj)
					// }
					if message.Action == "getConnection" {
						connectionObj := websocket.Connection{}
						if err := mapstructure.Decode(message.Message, &connectionObj); err != nil {
							log.Println(err)
						}
						connection <- connectionObj
					} else if message.Action == "pong" {
						// ignore
					} else {
						logObj := struct {
							CreatedAt string
							Resource  string
							Value     string
						}{}
						if err := mapstructure.Decode(message.Message, &logObj); err != nil {
							log.Println(err)
						} else {
							fmt.Printf("%s: %s %s\n", logObj.Resource, logObj.CreatedAt, logObj.Value)
						}
					}
				}
			}()

			message := websocket.Message{
				Action: "getConnection",
			}

			// We want to get the connectionId for further usage
			if err := c.WriteJSON(message); err != nil {
				panic(err)
			}
			ticker := time.NewTicker(time.Second * 10)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case connectionObj := <-connection:
					fmt.Printf("connected via connection %s\n\n", connectionObj.ID)
					fmt.Println("waiting for logs...")
					// Join function channels
					channels := []string{}
					for _, resource := range *Resources {
						channels = append(channels, fmt.Sprintf("log.%s", resource))
					}
					_, err := client.Client.Websocket.JoinChannels(&v1.JoinChannelsParams{
						ConnectionID: connectionObj.ID,
						Channels:     channels,
					})
					if err != nil {
						color.Red("err: %v", err)
						return
					}
				case <-ticker.C:
					if err := c.WriteJSON(websocket.Message{
						Action: "ping",
					}); err != nil {
						panic(err)
					}
				case <-interrupt:
					log.Println("interrupt")

					// Cleanly close the connection by sending a close message and then
					// waiting (with timeout) for the server to close the connection.
					err := c.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseNormalClosure, ""))
					if err != nil {
						log.Println("write close:", err)
						return
					}
					select {
					case <-done:
					case <-time.After(time.Second):
					}
					return
				}
			}
		},
	}
)

func init() {
	Resources = WatchCmd.Flags().StringSlice("resources", []string{}, `resources to watch logs for (e.g --resources "function:xyz")`)
}
