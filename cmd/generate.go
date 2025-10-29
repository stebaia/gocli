package cmd

import (
	"fmt"

	"fline-cli/internal/generator"
	"fline-cli/internal/ui"
	"fline-cli/internal/utils"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [feature-name]",
	Short: "Generate a feature (service, repository, BLoC)",
	Long: `Generate feature components following Pine architecture.

You can generate:
  • Service (Retrofit API client)
  • Repository (Data access layer)
  • BLoC (Business logic component)
  • All of the above

Example:
  pine generate user
  pine generate product --type service`,
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("name", "n", "", "Feature name")
	generateCmd.Flags().StringP("type", "t", "all", "Component type (service, repository, bloc, all)")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	logger := ui.NewLogger("generate")

	// Get feature name
	var featureName string
	if len(args) > 0 {
		featureName = args[0]
	} else {
		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			featureName = name
		}
	}

	// Interactive mode if no name provided
	if featureName == "" {
		var componentType string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Feature Name").
					Description("e.g., User, Product, Post").
					Value(&featureName).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("feature name is required")
						}
						return nil
					}),

				huh.NewSelect[string]().
					Title("Component Type").
					Description("What do you want to generate?").
					Options(
						huh.NewOption("All (Service + Repository + BLoC)", "all"),
						huh.NewOption("Service only", "service"),
						huh.NewOption("Repository only", "repository"),
						huh.NewOption("BLoC only", "bloc"),
					).
					Value(&componentType),
			).Title("🎨 Feature Generator"),
		)

		if err := form.Run(); err != nil {
			return err
		}

		cmd.Flags().Set("type", componentType)
	}

	componentType, _ := cmd.Flags().GetString("type")

	logger.Info(fmt.Sprintf("Generating %s for feature: %s", componentType, featureName))

	// Check if we're in a Flutter project
	writer := utils.NewFileWriter(".")
	if !writer.PathExists("pubspec.yaml") {
		logger.Error("Not in a Flutter project directory")
		logger.Info("Please run this command from your Flutter project root")
		return fmt.Errorf("pubspec.yaml not found")
	}

	// Get package name from pubspec
	packageName, err := getPackageName()
	if err != nil {
		return err
	}

	// Generate based on type
	gen := &FeatureGenerator{
		featureName: featureName,
		packageName: packageName,
		writer:      writer,
		logger:      logger,
	}

	switch componentType {
	case "service":
		if err := gen.generateService(); err != nil {
			return err
		}
	case "repository":
		if err := gen.generateRepository(); err != nil {
			return err
		}
	case "bloc":
		if err := gen.generateBloc(); err != nil {
			return err
		}
	case "all":
		if err := gen.generateAll(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid component type: %s", componentType)
	}

	logger.Success("Generation completed!")
	logger.Info("Don't forget to run: flutter pub run build_runner build --delete-conflicting-outputs")

	return nil
}

// FeatureGenerator generates feature components
type FeatureGenerator struct {
	featureName string
	packageName string
	writer      *utils.FileWriter
	logger      *ui.Logger
}

func (g *FeatureGenerator) generateAll() error {
	if err := g.generateService(); err != nil {
		return err
	}
	if err := g.generateRepository(); err != nil {
		return err
	}
	if err := g.generateBloc(); err != nil {
		return err
	}
	return nil
}

func (g *FeatureGenerator) generateService() error {
	naming := utils.NewNamingHelper(g.featureName)

	content := fmt.Sprintf(`import 'package:dio/dio.dart';
import 'package:retrofit/http.dart';

part '%s_service.g.dart';

@RestApi()
abstract class %sService {
  factory %sService(Dio dio) = _%sService;

  @GET('/endpoint')
  Future<List<dynamic>> get%ss();

  @GET('/endpoint/{id}')
  Future<dynamic> get%s(@Path('id') String id);

  @POST('/endpoint')
  Future<dynamic> create%s(@Body() Map<String, dynamic> data);

  @PUT('/endpoint/{id}')
  Future<dynamic> update%s(
    @Path('id') String id,
    @Body() Map<String, dynamic> data,
  );

  @DELETE('/endpoint/{id}')
  Future<void> delete%s(@Path('id') String id);
}
`,
		naming.SnakeCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
	)

	if err := g.writer.WriteFile(
		fmt.Sprintf("lib/network/service/%s_service.dart", naming.SnakeCase()),
		content,
	); err != nil {
		return err
	}

	g.logger.Success(fmt.Sprintf("Generated service: %s", naming.SnakeCase()))
	return nil
}

func (g *FeatureGenerator) generateRepository() error {
	naming := utils.NewNamingHelper(g.featureName)

	content := fmt.Sprintf(`import 'package:logger/logger.dart';
import 'package:%s/network/service/%s_service.dart';

class %sRepository {
  final %sService _service;
  final Logger _logger;

  %sRepository({
    required %sService service,
    required Logger logger,
  })  : _service = service,
        _logger = logger;

  Future<List<dynamic>> getAll() async {
    try {
      return await _service.get%ss();
    } catch (e) {
      _logger.e('Error fetching %ss', error: e);
      rethrow;
    }
  }

  Future<dynamic> getById(String id) async {
    try {
      return await _service.get%s(id);
    } catch (e) {
      _logger.e('Error fetching %s', error: e);
      rethrow;
    }
  }

  Future<dynamic> create(Map<String, dynamic> data) async {
    try {
      return await _service.create%s(data);
    } catch (e) {
      _logger.e('Error creating %s', error: e);
      rethrow;
    }
  }

  Future<dynamic> update(String id, Map<String, dynamic> data) async {
    try {
      return await _service.update%s(id, data);
    } catch (e) {
      _logger.e('Error updating %s', error: e);
      rethrow;
    }
  }

  Future<void> delete(String id) async {
    try {
      await _service.delete%s(id);
    } catch (e) {
      _logger.e('Error deleting %s', error: e);
      rethrow;
    }
  }
}
`,
		g.packageName,
		naming.SnakeCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
		naming.PascalCase(),
	)

	if err := g.writer.WriteFile(
		fmt.Sprintf("lib/repositories/%s_repository.dart", naming.SnakeCase()),
		content,
	); err != nil {
		return err
	}

	g.logger.Success(fmt.Sprintf("Generated repository: %s", naming.SnakeCase()))
	return nil
}

func (g *FeatureGenerator) generateBloc() error {
	naming := utils.NewNamingHelper(g.featureName)

	// Create a minimal model-less generator
	gen := generator.NewModelGenerator(
		g.featureName,
		map[string]interface{}{"id": "1"}, // Minimal JSON
		"/api/"+naming.KebabCase(),
		g.packageName,
		g.writer,
	)

	if err := gen.Generate(); err != nil {
		return err
	}

	g.logger.Success(fmt.Sprintf("Generated BLoC: %s", naming.SnakeCase()))
	return nil
}

func getPackageName() (string, error) {
	writer := utils.NewFileWriter(".")
	content, err := writer.ReadFile("pubspec.yaml")
	if err != nil {
		return "", err
	}

	// Simple parsing - look for "name: " line
	for _, line := range splitLines(content) {
		if len(line) > 6 && line[:5] == "name:" {
			return trimSpace(line[5:]), nil
		}
	}

	return "", fmt.Errorf("could not find package name in pubspec.yaml")
}

func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
