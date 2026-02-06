package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [url]",
	Short: "Connect to Sapliy Event Bus via WebSocket",
	Long:  `Connects to the Sapliy backend event bus to stream events in real-time.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverURL := "ws://localhost:8080/ws"
		if len(args) > 0 {
			serverURL = args[0]
		}

		apiKey, _ := cmd.Flags().GetString("key")
		trigger, _ := cmd.Flags().GetString("trigger")

		u, err := url.Parse(serverURL)
		if err != nil {
			log.Fatal("Invalid URL:", err)
		}

		fmt.Printf("🔌 Connecting to %s...\n", u.String())

		header := http.Header{}
		if apiKey != "" {
			header.Set("Authorization", "Bearer "+apiKey)
		}

		c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
			log.Fatal("Connection failed:", err)
		}
		defer c.Close()

		fmt.Println("✅ Connected! Listening for events...")

		done := make(chan struct{})

		// Reader loop
		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read-error:", err)
					return
				}
				fmt.Printf("< %s\n", message)
			}
		}()

		// Trigger logic
		if trigger != "" {
			fmt.Printf("> Triggering event: %s\n", trigger)
			err := c.WriteMessage(websocket.TextMessage, []byte(trigger))
			if err != nil {
				log.Println("write-error:", err)
			}
		}

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		for {
			select {
			case <-done:
				return
			case <-interrupt:
				fmt.Println("\nDisconnecting...")
				// Cleanly close the connection by sending a close message
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write-close:", err)
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

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringP("key", "k", "", "API Key for authentication")
	connectCmd.Flags().StringP("trigger", "t", "", "Send a JSON event payload immediately after connecting")
}
