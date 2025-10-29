package cmd

import (
	"fmt"
	"os"

	"fline-cli/internal/ui"
	"fline-cli/internal/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pine",
	Short: "Pine CLI - Flutter Project Generator",
	Long: `Pine CLI is a powerful Flutter project generator with Pine architecture.

It helps you create Flutter projects with:
  • Clean Architecture (Service, Repository, BLoC, UI)
  • Firebase or Supabase integration
  • Push notifications setup
  • Auto-generated models from JSON
  • Beautiful example screens

Made with ❤️ for Flutter developers`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check if Flutter is installed
		if err := utils.CheckFlutterInstalled(); err != nil {
			logger := ui.NewLogger("pine")
			logger.Error("Flutter is not installed or not in PATH")
			logger.Info("Please install Flutter: https://flutter.dev/docs/get-started/install")
			os.Exit(1)
		}
	},
}

func Execute() {
	// Print banner
	ui.PrintBanner()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add version info
	rootCmd.Version = "1.0.0"
	rootCmd.SetVersionTemplate(fmt.Sprintf("Pine CLI v%s\n", rootCmd.Version))
}
