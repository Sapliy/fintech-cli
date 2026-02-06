package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Sapliy resources (zones, flows)",
	Long:  `Scaffold configuration files for Sapliy Automation Zones and Flows.`,
}

var zoneCmd = &cobra.Command{
	Use:   "zone [name]",
	Short: "Generate a new automation zone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fileName := fmt.Sprintf("%s.zone.json", strings.ToLower(name))

		content := fmt.Sprintf(`{
  "id": "zone_%s",
  "name": "%s",
  "description": "Automation zone for %s",
  "version": "1.0.0",
  "triggers": [],
  "actions": []
}`, name, name, name)

		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			fmt.Printf("Error creating zone: %v\n", err)
			return
		}
		fmt.Printf("✅ Generated zone file: %s\n", fileName)
	},
}

var flowCmd = &cobra.Command{
	Use:   "flow [name]",
	Short: "Generate a new automation flow",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fileName := fmt.Sprintf("%s.flow.json", strings.ToLower(name))

		content := fmt.Sprintf(`{
  "id": "flow_%s",
  "name": "%s",
  "steps": [
    {
      "id": "start",
      "type": "trigger",
      "config": {}
    }
  ]
}`, name, name)

		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			fmt.Printf("Error creating flow: %v\n", err)
			return
		}
		fmt.Printf("✅ Generated flow file: %s\n", fileName)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(zoneCmd)
	generateCmd.AddCommand(flowCmd)
}
