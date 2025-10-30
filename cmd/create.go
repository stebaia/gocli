package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"fline-cli/internal/config"
	"fline-cli/internal/generator"
	"fline-cli/internal/ui"
	"fline-cli/internal/utils"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Flutter project with Pine architecture",
	Long: `Create a new Flutter project with Pine architecture.

This interactive wizard will guide you through:
  • Project configuration
  • Backend selection (Firebase/Supabase)
  • Notification setup
  • Model generation from JSON
  • Example screens generation`,
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Flags for non-interactive mode
	createCmd.Flags().StringP("name", "n", "", "Project name")
	createCmd.Flags().StringP("path", "p", ".", "Target directory path")
	createCmd.Flags().StringP("org", "o", "", "Organization name (e.g., com.example)")
	createCmd.Flags().BoolP("force", "f", false, "Force creation even if directory exists")
	createCmd.Flags().Bool("firebase", false, "Enable Firebase integration")
	createCmd.Flags().Bool("supabase", false, "Enable Supabase integration")
	createCmd.Flags().Bool("no-interactive", false, "Disable interactive mode")
}

func runCreate(cmd *cobra.Command, args []string) error {
	logger := ui.NewLogger("create")

	// Check if non-interactive mode
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")

	var cfg *config.ProjectConfig
	var err error

	if noInteractive {
		cfg, err = createFromFlags(cmd)
	} else {
		cfg, err = runInteractiveWizard()
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Configuration error: %s", err))
		return err
	}

	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		logger.Error(fmt.Sprintf("Validation error: %s", err))
		return err
	}

	// Show summary
	showProjectSummary(cfg)

	// Confirm
	if !noInteractive {
		var confirm bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Ready to create your project?").
					Value(&confirm),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		if !confirm {
			logger.Warning("Project creation cancelled")
			return nil
		}
	}

	// Generate project
	logger.Title("Creating your Flutter project...")

	gen := generator.NewProjectGenerator(cfg)
	if err := gen.Generate(); err != nil {
		logger.Error(fmt.Sprintf("Failed to generate project: %s", err))
		return err
	}

	return nil
}

func runInteractiveWizard() (*config.ProjectConfig, error) {
	cfg := config.DefaultProjectConfig()

	// Basic Information
	basicForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Description("Use lowercase with underscores (e.g., my_app)").
				Value(&cfg.ProjectName).
				Validate(func(s string) error {
					return utils.ValidateProjectName(s)
				}),

			huh.NewInput().
				Title("Target Directory").
				Description("Where to create the project (default: current directory)").
				Value(&cfg.TargetDirectory).
				Placeholder("."),

			huh.NewInput().
				Title("Organization").
				Description("Reverse domain notation (e.g., com.example)").
				Value(&cfg.OrganizationName).
				Validate(func(s string) error {
					return utils.ValidateOrganization(s)
				}),

			huh.NewInput().
				Title("Description").
				Description("A short description of your project").
				Value(&cfg.Description),
		).Title("📱 Basic Information").
			Description("Let's start with the basics"),
	)

	if err := basicForm.Run(); err != nil {
		return nil, err
	}

	// Backend Selection
	var backendChoice string
	backendForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Backend Service").
				Description("Choose your backend service").
				Options(
					huh.NewOption("None", "none"),
					huh.NewOption("Firebase", "firebase"),
					huh.NewOption("Supabase", "supabase"),
				).
				Value(&backendChoice),
		).Title("🔥 Backend Configuration"),
	)

	if err := backendForm.Run(); err != nil {
		return nil, err
	}

	cfg.UseFirebase = backendChoice == "firebase"
	cfg.UseSupabase = backendChoice == "supabase"

	// Notifications Setup
	if cfg.UseFirebase {
		notifForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Enable Push Notifications?").
					Description("Setup Firebase Cloud Messaging").
					Value(&cfg.EnableNotifications),
			).Title("🔔 Notifications"),
		)

		if err := notifForm.Run(); err != nil {
			return nil, err
		}

		if cfg.EnableNotifications {
			cfg.NotificationService = "fcm"
		}
	}

	// Models from JSON
	var addModels bool
	modelsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add models from JSON?").
				Description("Generate models, services, repositories and BLoCs from JSON data").
				Value(&addModels),
		).Title("📦 Models"),
	)

	if err := modelsForm.Run(); err != nil {
		return nil, err
	}

	if addModels {
		if err := addModelsInteractive(cfg); err != nil {
			return nil, err
		}
	}

	// Example Screens
	var selectedScreens []string
	screensForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Example Screens").
				Description("Use SPACE to select/unselect, ENTER to confirm").
				Options(
					huh.NewOption("Login Screen", "login"),
					huh.NewOption("Home Screen", "home"),
					huh.NewOption("Profile Screen", "profile"),
					huh.NewOption("Settings Screen", "settings"),
				).
				Value(&selectedScreens),
		).Title("🎨 Screens"),
	)

	if err := screensForm.Run(); err != nil {
		return nil, err
	}

	// Set screen flags
	for _, screen := range selectedScreens {
		switch screen {
		case "login":
			cfg.GenerateLoginScreen = true
		case "home":
			cfg.GenerateHomeScreen = true
		case "profile":
			cfg.GenerateProfileScreen = true
		case "settings":
			cfg.GenerateSettingsScreen = true
		}
	}

	return cfg, nil
}

func addModelsInteractive(cfg *config.ProjectConfig) error {
	for {
		var modelName, jsonStr, endpoint string
		var addAnother bool

		modelForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Model Name").
					Description("e.g., User, Product, Post").
					Value(&modelName),

				huh.NewText().
					Title("JSON Data").
					Description("Paste your JSON object here").
					Value(&jsonStr).
					Lines(10),

				huh.NewInput().
					Title("API Endpoint").
					Description("e.g., /api/users").
					Value(&endpoint),

				huh.NewConfirm().
					Title("Add another model?").
					Value(&addAnother),
			).Title(fmt.Sprintf("📝 Model %d", len(cfg.Models)+1)),
		)

		if err := modelForm.Run(); err != nil {
			return err
		}

		// Parse JSON
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
			fmt.Println(ui.ErrorStyle.Render("Invalid JSON format. Please try again."))
			continue
		}

		cfg.Models = append(cfg.Models, config.ModelConfig{
			Name:     modelName,
			JSONData: jsonData,
			Endpoint: endpoint,
		})

		if !addAnother {
			break
		}
	}

	return nil
}

func createFromFlags(cmd *cobra.Command) (*config.ProjectConfig, error) {
	cfg := config.DefaultProjectConfig()

	name, _ := cmd.Flags().GetString("name")
	targetPath, _ := cmd.Flags().GetString("path")
	org, _ := cmd.Flags().GetString("org")
	force, _ := cmd.Flags().GetBool("force")
	firebase, _ := cmd.Flags().GetBool("firebase")
	supabase, _ := cmd.Flags().GetBool("supabase")

	if name == "" {
		return nil, fmt.Errorf("project name is required (use --name flag)")
	}

	cfg.ProjectName = name
	cfg.TargetDirectory = targetPath
	if org != "" {
		cfg.OrganizationName = org
	}
	cfg.Force = force
	cfg.UseFirebase = firebase
	cfg.UseSupabase = supabase

	return cfg, nil
}

func validateConfig(cfg *config.ProjectConfig) error {
	if err := utils.ValidateProjectName(cfg.ProjectName); err != nil {
		return err
	}

	if err := utils.ValidateOrganization(cfg.OrganizationName); err != nil {
		return err
	}

	// Set default target directory
	if cfg.TargetDirectory == "" {
		cfg.TargetDirectory = "."
	}

	// Check if target directory exists
	if _, err := os.Stat(cfg.TargetDirectory); os.IsNotExist(err) {
		return fmt.Errorf("target directory '%s' does not exist", cfg.TargetDirectory)
	}

	// Build full project path
	projectPath := cfg.TargetDirectory
	if projectPath != "." {
		projectPath = projectPath + "/" + cfg.ProjectName
	} else {
		projectPath = cfg.ProjectName
	}

	// Check if project directory exists and force is not set
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) && !cfg.Force {
		return fmt.Errorf("directory '%s' already exists (use --force to overwrite)", projectPath)
	}

	return nil
}

func showProjectSummary(cfg *config.ProjectConfig) {
	logger := ui.NewLogger("summary")

	logger.Title("Project Summary")

	// Build project path for display
	projectPath := cfg.TargetDirectory
	if projectPath != "." {
		projectPath = projectPath + "/" + cfg.ProjectName
	} else {
		projectPath = cfg.ProjectName
	}

	items := []string{
		fmt.Sprintf("Name: %s", cfg.ProjectName),
		fmt.Sprintf("Location: %s", projectPath),
		fmt.Sprintf("Organization: %s", cfg.OrganizationName),
	}

	if cfg.Description != "" {
		items = append(items, fmt.Sprintf("Description: %s", cfg.Description))
	}

	if cfg.UseFirebase {
		items = append(items, ui.SuccessStyle.Render("✓")+" Firebase integration")
	}

	if cfg.UseSupabase {
		items = append(items, ui.SuccessStyle.Render("✓")+" Supabase integration")
	}

	if cfg.EnableNotifications {
		items = append(items, ui.SuccessStyle.Render("✓")+" Push Notifications ("+cfg.NotificationService+")")
	}

	if len(cfg.Models) > 0 {
		modelNames := []string{}
		for _, m := range cfg.Models {
			modelNames = append(modelNames, m.Name)
		}
		items = append(items, fmt.Sprintf("Models: %s", strings.Join(modelNames, ", ")))
	}

	screens := []string{}
	if cfg.GenerateLoginScreen {
		screens = append(screens, "Login")
	}
	if cfg.GenerateHomeScreen {
		screens = append(screens, "Home")
	}
	if cfg.GenerateProfileScreen {
		screens = append(screens, "Profile")
	}
	if cfg.GenerateSettingsScreen {
		screens = append(screens, "Settings")
	}

	if len(screens) > 0 {
		items = append(items, fmt.Sprintf("Screens: %s", strings.Join(screens, ", ")))
	}

	logger.Box("Configuration", items)
}
