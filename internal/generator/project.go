package generator

import (
	"fmt"

	"fline-cli/internal/config"
	"fline-cli/internal/templates"
	"fline-cli/internal/ui"
	"fline-cli/internal/utils"
)

// ProjectGenerator generates a complete Flutter project
type ProjectGenerator struct {
	config *config.ProjectConfig
	logger *ui.Logger
	writer *utils.FileWriter
	flutter *utils.FlutterCLI
}

// NewProjectGenerator creates a new project generator
func NewProjectGenerator(cfg *config.ProjectConfig) *ProjectGenerator {
	return &ProjectGenerator{
		config: cfg,
		logger: ui.NewLogger("generator"),
	}
}

// Generate generates the complete project
func (g *ProjectGenerator) Generate() error {
	g.logger.Step(1, 8, "Creating Flutter project...")
	if err := g.createFlutterProject(); err != nil {
		return fmt.Errorf("failed to create Flutter project: %w", err)
	}

	// Build project path
	projectPath := g.config.TargetDirectory
	if projectPath != "." {
		projectPath = projectPath + "/" + g.config.ProjectName
	} else {
		projectPath = g.config.ProjectName
	}

	// Initialize helpers
	g.writer = utils.NewFileWriter(projectPath)
	g.flutter = utils.NewFlutterCLI(projectPath)

	g.logger.Step(2, 8, "Updating pubspec.yaml...")
	if err := g.updatePubspec(); err != nil {
		return fmt.Errorf("failed to update pubspec: %w", err)
	}

	g.logger.Step(3, 8, "Creating folder structure...")
	if err := g.createFolderStructure(); err != nil {
		return fmt.Errorf("failed to create folders: %w", err)
	}

	g.logger.Step(4, 8, "Generating core files...")
	if err := g.generateCoreFiles(); err != nil {
		return fmt.Errorf("failed to generate core files: %w", err)
	}

	g.logger.Step(5, 8, "Installing dependencies...")
	if err := g.flutter.PubGet(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	// Firebase integration
	if g.config.UseFirebase {
		g.logger.Step(6, 8, "Setting up Firebase...")
		if err := g.setupFirebase(); err != nil {
			return fmt.Errorf("failed to setup Firebase: %w", err)
		}
	}

	// Supabase integration
	if g.config.UseSupabase {
		g.logger.Step(6, 8, "Setting up Supabase...")
		if err := g.setupSupabase(); err != nil {
			return fmt.Errorf("failed to setup Supabase: %w", err)
		}
	}

	// Generate models
	if len(g.config.Models) > 0 {
		g.logger.Step(7, 8, fmt.Sprintf("Generating %d models...", len(g.config.Models)))
		if err := g.generateModels(); err != nil {
			return fmt.Errorf("failed to generate models: %w", err)
		}
	}

	// Generate screens
	g.logger.Step(7, 8, "Generating screens...")
	if err := g.generateScreens(); err != nil {
		return fmt.Errorf("failed to generate screens: %w", err)
	}

	g.logger.Step(8, 8, "Running code generation...")
	if err := g.runCodeGeneration(); err != nil {
		g.logger.Warning("Code generation had issues, but you can run it manually later")
		g.logger.Info("Run: flutter pub run build_runner build --delete-conflicting-outputs")
	}

	return nil
}

func (g *ProjectGenerator) createFlutterProject() error {
	flutter := utils.NewFlutterCLI(g.config.TargetDirectory)
	return flutter.Create(
		g.config.ProjectName,
		g.config.OrganizationName,
		g.config.Force,
	)
}

func (g *ProjectGenerator) updatePubspec() error {
	content := templates.GeneratePubspec(g.config)
	return g.writer.WriteFile("pubspec.yaml", content)
}

func (g *ProjectGenerator) createFolderStructure() error {
	folders := []string{
		"lib/di",
		"lib/l10n",
		"lib/mappers",
		"lib/model",
		"lib/network/interceptor",
		"lib/network/service",
		"lib/repositories",
		"lib/routers",
		"lib/state_management/bloc",
		"lib/state_management/cubit",
		"lib/state_management/provider",
		"lib/ui",
		"lib/utils",
		"lib/theme",
	}

	for _, folder := range folders {
		if err := g.writer.EnsureDir(folder); err != nil {
			return err
		}
	}

	return nil
}

func (g *ProjectGenerator) generateCoreFiles() error {
	files := map[string]string{
		"lib/main.dart":                       templates.GenerateMain(),
		"lib/app.dart":                        templates.GenerateApp(g.config.ProjectName),
		"lib/di/dependency_injector.dart":     templates.GenerateDependencyInjector(),
		"lib/di/blocs.dart":                   templates.GenerateBlocs(),
		"lib/di/mappers.dart":                 templates.GenerateMappers(),
		"lib/di/providers.dart":               templates.GenerateProviders(),
		"lib/di/repositories.dart":            templates.GenerateRepositories(),
		"lib/theme/light_theme.dart":          templates.GenerateLightTheme(),
		"lib/routers/app_router.dart":         templates.GenerateAppRouter(g.config.ProjectName),
		"lib/l10n/app_en.arb":                 templates.GenerateL10n("en", g.config.ProjectName),
		"lib/l10n/app_it.arb":                 templates.GenerateL10n("it", g.config.ProjectName),
	}

	for path, content := range files {
		if err := g.writer.WriteFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func (g *ProjectGenerator) setupFirebase() error {
	gen := NewFirebaseGenerator(g.config, g.writer, g.flutter)
	return gen.Generate()
}

func (g *ProjectGenerator) setupSupabase() error {
	gen := NewSupabaseGenerator(g.config, g.writer, g.flutter)
	return gen.Generate()
}

func (g *ProjectGenerator) generateModels() error {
	for _, model := range g.config.Models {
		gen := NewModelGenerator(
			model.Name,
			model.JSONData,
			model.Endpoint,
			g.config.ProjectName,
			g.writer,
		)

		if err := gen.Generate(); err != nil {
			return fmt.Errorf("failed to generate model %s: %w", model.Name, err)
		}

		g.logger.Success(fmt.Sprintf("Generated model: %s", model.Name))
	}

	return nil
}

func (g *ProjectGenerator) generateScreens() error {
	gen := NewScreenGenerator(g.config, g.writer)
	return gen.Generate()
}

func (g *ProjectGenerator) runCodeGeneration() error {
	// Generate localizations
	if err := g.flutter.GenL10n(); err != nil {
		return err
	}

	// Run build_runner
	if err := g.flutter.BuildRunnerBuild(); err != nil {
		return err
	}

	return nil
}
