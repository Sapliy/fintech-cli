package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// Embedded static files from Next.js export
// We assume the user creates an 'out' directory in sapliy-cli/pkg/cmd/ui via their build script
// Since we don't have the build artifact yet, we use a placeholder variable or filesystem fallback.
// Use 'embed' for single binary distribution.
//
//go:embed ui/*
var content embed.FS

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Sapliy Automation Studio locally",
	Long:  `Hosts the self-contained Sapliy Automation Studio web interface locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")

		// Setup filesystem
		// If embedded 'ui' folder exists and has content, use it.
		// Otherwise, if running locally and 'ui' is empty/missing, maybe serve from relative path?
		// For robust embedding, we strip the prefix.

		fsys, err := fs.Sub(content, "ui")
		if err != nil {
			log.Fatal(err)
		}

		// Fallback for dev: check if we are in source and have ../fintech-automation/out?
		// But embedding is for distribution.

		http.Handle("/", http.FileServer(http.FS(fsys)))

		fmt.Printf("🚀 Sapliy Automation Studio running at http://localhost:%s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("port", "p", "3000", "Port to serve the studio on")
}
