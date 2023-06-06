package builder

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/heeser-io/universe-cli/config"
)

func (b *Builder) Logs() {
	apiKey := config.Main.GetString("apiKey")

	u := url.URL{
		Scheme:   "wss",
		Host:     "https://api.universecloud.dev/websocket",
		Path:     "/websocket",
		RawQuery: fmt.Sprintf("devkey=%s&subscribeTo=*.*.log.create.function:62ef7a882ec831eecae22100", apiKey),
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	done := make(chan struct{})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go func() {
		// receive messages
		// for {
		// 	msg := v1.IncomingMessage{}
		// 	if err := c.ReadJSON(&msg); err != nil {
		// 		log.Println("read: ", err)
		// 		os.Exit(1)
		// 		return
		// 	}

		// 	log := v1.Log{}

		// 	if err := json.Unmarshal([]byte(msg.Message), &log); err != nil {
		// 		fmt.Printf("err: %v\n", err)
		// 	}

		// 	fmt.Printf("%s: %s\n", log.CreatedAt, log.Value.(string))
		// }
	}()
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
}
