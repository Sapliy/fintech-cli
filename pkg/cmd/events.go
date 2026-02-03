package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/sapliy/fintech-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var eventData string
var zoneID string

var triggerCmd = &cobra.Command{
	Use:   "trigger [event_type]",
	Short: "Trigger a mock event for automation flows",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			fmt.Println("Error: API key not set. Use 'sapliy auth login' or set in config.")
			return
		}

		eventType := args[0]

		var data map[string]interface{}
		if eventData != "" {
			if err := json.Unmarshal([]byte(eventData), &data); err != nil {
				log.Fatalf("Invalid JSON data: %v", err)
			}
		}

		client := fintech.NewClient(apiKey, fintech.WithBaseURL(viper.GetString("api_url")))

		// In a real implementation, this would hit a dedicated trigger endpoint
		// For now, we'll simulate the call
		fmt.Printf("Triggering event '%s' in zone '%s'...\n", eventType, zoneID)

		// Use the new SDK TriggerEvent method
		err := client.TriggerEvent(context.Background(), eventType, zoneID, data)

		if err != nil {
			fmt.Printf("Failed to trigger event: %v\n", err)
			return
		}

		fmt.Println("✅ Event triggered successfully! The Flow Runner will process it shortly.")
	},
}

func init() {
	rootCmd.AddCommand(triggerCmd)
	triggerCmd.Flags().StringVarP(&eventData, "data", "d", "{}", "JSON event data")
	triggerCmd.Flags().StringVarP(&zoneID, "zone", "z", "", "Zone ID to scope the event")
	triggerCmd.MarkFlagRequired("zone")
}
