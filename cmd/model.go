package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"fline-cli/internal/generator"
	"fline-cli/internal/ui"
	"fline-cli/internal/utils"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:   "model [name]",
	Short: "Generate model, service, repository, and BLoC from JSON",
	Long: `Generate a complete feature from JSON data.

This command creates:
  • Model with json_annotation
  • Retrofit service
  • Repository with error handling
  • BLoC with all CRUD operations

Example:
  pine model user --json '{"id": 1, "name": "John", "email": "john@example.com"}'
  pine model product --json-file product.json`,
	RunE: runModel,
}

func init() {
	rootCmd.AddCommand(modelCmd)

	modelCmd.Flags().StringP("name", "n", "", "Model name")
	modelCmd.Flags().StringP("json", "j", "", "JSON string")
	modelCmd.Flags().StringP("json-file", "f", "", "Path to JSON file")
	modelCmd.Flags().StringP("endpoint", "e", "", "API endpoint (e.g., /api/users)")
}

func runModel(cmd *cobra.Command, args []string) error {
	logger := ui.NewLogger("model")

	// Get model name
	var modelName string
	if len(args) > 0 {
		modelName = args[0]
	} else {
		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			modelName = name
		}
	}

	// Get JSON data
	jsonStr, _ := cmd.Flags().GetString("json")
	jsonFile, _ := cmd.Flags().GetString("json-file")
	endpoint, _ := cmd.Flags().GetString("endpoint")

	// Interactive mode if parameters missing
	if modelName == "" || (jsonStr == "" && jsonFile == "") {
		var err error
		modelName, jsonStr, endpoint, err = runModelInteractive()
		if err != nil {
			return err
		}
	}

	// Check if we're in a Flutter project
	writer := utils.NewFileWriter(".")
	if !writer.PathExists("pubspec.yaml") {
		logger.Error("Not in a Flutter project directory")
		logger.Info("Please run this command from your Flutter project root")
		return fmt.Errorf("pubspec.yaml not found")
	}

	// Get package name
	packageName, err := getPackageName()
	if err != nil {
		return err
	}

	// Parse JSON
	var jsonData map[string]interface{}
	if jsonFile != "" {
		// Read from file
		content, err := os.ReadFile(jsonFile)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to read JSON file: %s", err))
			return err
		}
		jsonStr = string(content)
	}

	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		logger.Error("Invalid JSON format")
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Set default endpoint if not provided
	if endpoint == "" {
		naming := utils.NewNamingHelper(modelName)
		endpoint = "/api/" + naming.KebabCase()
	}

	logger.Info(fmt.Sprintf("Generating model: %s", modelName))
	logger.Info(fmt.Sprintf("Endpoint: %s", endpoint))

	// Generate
	gen := generator.NewModelGenerator(
		modelName,
		jsonData,
		endpoint,
		packageName,
		writer,
	)

	if err := gen.Generate(); err != nil {
		logger.Error(fmt.Sprintf("Generation failed: %s", err))
		return err
	}

	logger.Success("Model generated successfully!")
	logger.NewLine()
	logger.Box("Generated files:", []string{
		fmt.Sprintf("lib/model/%s.dart", utils.NewNamingHelper(modelName).SnakeCase()),
		fmt.Sprintf("lib/network/service/%s_service.dart", utils.NewNamingHelper(modelName).SnakeCase()),
		fmt.Sprintf("lib/repositories/%s_repository.dart", utils.NewNamingHelper(modelName).SnakeCase()),
		fmt.Sprintf("lib/state_management/bloc/%s/*", utils.NewNamingHelper(modelName).SnakeCase()),
	})
	logger.NewLine()
	logger.Info("Next steps:")
	logger.Info("1. Run: flutter pub run build_runner build --delete-conflicting-outputs")
	logger.Info("2. Add the service, repository, and BLoC to your dependency injector")

	return nil
}

func runModelInteractive() (string, string, string, error) {
	var modelName, jsonStr, endpoint string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Model Name").
				Description("e.g., User, Product, Post").
				Value(&modelName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("model name is required")
					}
					return nil
				}),

			huh.NewText().
				Title("JSON Data").
				Description("Paste your JSON object here").
				Value(&jsonStr).
				Lines(10).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("JSON data is required")
					}
					var tmp map[string]interface{}
					if err := json.Unmarshal([]byte(s), &tmp); err != nil {
						return fmt.Errorf("invalid JSON format")
					}
					return nil
				}),

			huh.NewInput().
				Title("API Endpoint").
				Description("e.g., /api/users (optional, will auto-generate if empty)").
				Value(&endpoint),
		).Title("📦 Model Generator").
			Description("Generate a complete feature from JSON"),
	)

	if err := form.Run(); err != nil {
		return "", "", "", err
	}

	return modelName, jsonStr, endpoint, nil
}
