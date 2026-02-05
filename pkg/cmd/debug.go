package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug and inspect flows in real-time",
	Long: `Debug command provides real-time inspection of flows and events.
Connect to the Sapliy event stream to monitor automation flows as they execute.`,
}

var debugListenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen to real-time event stream via WebSocket",
	Long: `Connect to Sapliy API and stream events in real-time.
This is useful for debugging flows and watching events as they happen.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			fmt.Println("Error: API key not set. Use 'sapliy auth login'.")
			os.Exit(1)
		}

		zone := viper.GetString("current_zone")
		if zone == "" {
			zone, _ = cmd.Flags().GetString("zone")
		}

		// Determine WS URL (default to localhost:8089 for dev)
		apiURL := viper.GetString("api_url")
		wsURL := "ws://localhost:8089/v1/events/stream"
		if apiURL != "" && !strings.Contains(apiURL, "localhost") {
			// Production logic would replace https:// with wss://
			wsURL = strings.Replace(apiURL, "https://", "wss://", 1) + "/v1/events/stream"
		}

		// Append query params
		wsURL += fmt.Sprintf("?api_key=%s", apiKey)
		if zone != "" {
			wsURL += fmt.Sprintf("&zone=%s", zone)
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		filterType, _ := cmd.Flags().GetString("filter")

		fmt.Printf("🔌 Connecting to %s...\n", wsURL)

		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			fmt.Printf("❌ Failed to connect: %v\n", err)
			return
		}
		defer conn.Close()

		fmt.Println("✅ Connected! Streaming events... (Ctrl+C to stop)")
		fmt.Println(strings.Repeat("─", 60))

		// Handle graceful shutdown
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					// Check if normal close
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						fmt.Printf("❌ connection error: %v\n", err)
					}
					return
				}

				var event map[string]interface{}
				if err := json.Unmarshal(message, &event); err != nil {
					continue
				}

				eventType, _ := event["type"].(string)

				// Apply filter if specified
				if filterType != "" && !strings.Contains(eventType, filterType) {
					continue
				}

				timestamp := time.Now().Format("15:04:05")

				if verbose {
					prettyJSON, _ := json.MarshalIndent(event, "", "  ")
					fmt.Printf("[%s] %s\n%s\n\n", timestamp, eventType, string(prettyJSON))
				} else {
					// Try to get ID if available
					id := ""
					if data, ok := event["data"].(map[string]interface{}); ok {
						if val, ok := data["id"].(string); ok {
							id = val
						}
					}
					fmt.Printf("[%s] %-30s  %s\n", timestamp, eventType, id)
				}
			}
		}()

		select {
		case <-interrupt:
			fmt.Println("\n👋 Disconnecting...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
		case <-done:
			fmt.Println("Server closed connection")
		}
	},
}

// pollEvents fetches events from the API

var debugInspectCmd = &cobra.Command{
	Use:   "inspect [flow_id]",
	Short: "Inspect a specific flow execution",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			fmt.Println("Error: API key not set.")
			os.Exit(1)
		}

		flowID := args[0]
		fmt.Printf("🔍 Inspecting flow: %s\n", flowID)
		fmt.Println(strings.Repeat("─", 60))

		// TODO: Implement API call to get flow details
		fmt.Println("Flow inspection coming soon...")
	},
}

var debugReplCmd = &cobra.Command{
	Use:   "repl",
	Short: "Interactive REPL for testing events",
	Long: `Start an interactive REPL to test events and flows.
Type event types and JSON data to trigger events interactively.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			fmt.Println("Error: API key not set.")
			os.Exit(1)
		}

		zone := viper.GetString("current_zone")

		fmt.Println("🎮 Sapliy Debug REPL")
		fmt.Println("Type 'help' for commands, 'exit' to quit")
		fmt.Printf("Current zone: %s\n", zone)
		fmt.Println(strings.Repeat("─", 60))

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("sapliy> ")
			if !scanner.Scan() {
				break
			}

			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}

			switch input {
			case "exit", "quit":
				fmt.Println("👋 Goodbye!")
				return
			case "help":
				fmt.Println(`Commands:
  emit <type> [json]  - Emit an event (e.g., emit payment.created {"amount":100})
  zone <id>           - Switch to a different zone
  status              - Show current configuration
  exit                - Exit the REPL`)
			case "status":
				fmt.Printf("API Key: %s...%s\n", apiKey[:8], apiKey[len(apiKey)-4:])
				fmt.Printf("Zone: %s\n", zone)
				fmt.Printf("API URL: %s\n", viper.GetString("api_url"))
			default:
				if strings.HasPrefix(input, "emit ") {
					parts := strings.SplitN(input[5:], " ", 2)
					eventType := parts[0]
					data := "{}"
					if len(parts) > 1 {
						data = parts[1]
					}
					fmt.Printf("➡️  Emitting %s: %s\n", eventType, data)
					// TODO: Actually emit the event via SDK
				} else if strings.HasPrefix(input, "zone ") {
					zone = strings.TrimSpace(input[5:])
					viper.Set("current_zone", zone)
					fmt.Printf("✅ Switched to zone: %s\n", zone)
				} else {
					fmt.Printf("Unknown command: %s\n", input)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.AddCommand(debugListenCmd)
	debugCmd.AddCommand(debugInspectCmd)
	debugCmd.AddCommand(debugReplCmd)

	debugListenCmd.Flags().StringP("zone", "z", "", "Zone ID to filter events")
	debugListenCmd.Flags().BoolP("verbose", "v", false, "Show full event payloads")
	debugListenCmd.Flags().StringP("filter", "f", "", "Filter events by type (substring match)")
}
