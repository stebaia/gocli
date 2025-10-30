package templates

import (
	"fline-cli/internal/config"
	"fmt"
	"strings"
)

// GeneratePubspec generates pubspec.yaml content
func GeneratePubspec(cfg *config.ProjectConfig) string {
	dependencies := []string{
		"flutter:",
		"  sdk: flutter",
		"flutter_localizations:",
		"  sdk: flutter",
		"cupertino_icons: ^1.0.8",
		"logger: ^2.5.0",
		"flutter_secure_storage: ^9.2.2",
		"flutter_bloc: ^8.1.6",
		"hydrated_bloc: ^9.1.5",
		"equatable: ^2.0.7",
		"font_awesome_flutter: ^10.8.0",
		"pine: ^1.0.3",
		"provider: ^6.0.5",
		"retrofit: ^4.4.1",
		"dio: ^5.7.0",
		"pretty_dio_logger: ^1.4.0",
		"auto_route: ^9.2.2",
		"cached_network_image: ^3.2.3",
		"json_annotation: ^4.9.0",
		"sqlite3_flutter_libs: ^0.5.27",
		"path_provider: ^2.1.5",
		"path: ^1.9.0",
		"shared_preferences: ^2.3.3",
		"intl: ^0.20.2",
		"google_fonts: ^6.3.2",
		"flutter_local_notifications: ^18.0.1",
	}

	// Add Firebase dependencies
	if cfg.UseFirebase {
		dependencies = append(dependencies,
			"firebase_core: ^4.1.1",
			"firebase_auth: ^6.1.0",
			"cloud_firestore: ^6.0.2",
			"google_sign_in: ^7.2.0",
		)

		if cfg.EnableNotifications {
			dependencies = append(dependencies, "firebase_messaging: ^15.1.6")
		}
	}

	// Add Supabase dependencies
	if cfg.UseSupabase {
		dependencies = append(dependencies, "supabase_flutter: ^2.9.1")
	}

	devDependencies := []string{
		"flutter_test:",
		"  sdk: flutter",
		"flutter_lints: ^5.0.0",
		"build_runner: ^2.4.13",
		"bloc_test: ^9.1.0",
		"retrofit_generator: ^9.1.5",
		"auto_route_generator: ^9.0.0",
		"http_mock_adapter: ^0.6.1",
		"data_fixture_dart: ^2.2.0",
		"mockito: ^5.3.2",
		"json_serializable: ^6.7.1",
	}

	description := cfg.Description
	if description == "" {
		description = "A new Flutter project with Pine architecture."
	}

	return fmt.Sprintf(`name: %s
description: %s
publish_to: 'none'
version: 1.0.0+1

environment:
  sdk: '>=3.0.0 <4.0.0'

dependencies:
  %s

dev_dependencies:
  %s

flutter:
  uses-material-design: true
  generate: true
`,
		cfg.ProjectName,
		description,
		strings.Join(dependencies, "\n  "),
		strings.Join(devDependencies, "\n  "),
	)
}

// GenerateMain generates main.dart
func GenerateMain() string {
	return `import 'package:flutter/material.dart';
import 'app.dart';

void main() {
  runApp(const App());
}
`
}

// GenerateApp generates app.dart
func GenerateApp(packageName string) string {
	return fmt.Sprintf(`import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:%s/l10n/app_localizations.dart';
import 'package:%s/di/dependency_injector.dart';
import 'package:%s/routers/app_router.dart';
import 'package:%s/theme/light_theme.dart';

final router = AppRouter();

class App extends StatelessWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context) {
    SystemChrome.setPreferredOrientations([
      DeviceOrientation.portraitUp,
      DeviceOrientation.portraitDown,
    ]);

    return DependencyInjector(
      child: MaterialApp.router(
        debugShowCheckedModeBanner: false,
        routeInformationParser: router.defaultRouteParser(),
        routerDelegate: router.delegate(),
        localizationsDelegates: const [
          AppLocalizations.delegate,
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        theme: LightTheme.make,
        supportedLocales: const [
          Locale('en'),
          Locale('it'),
        ],
      ),
    );
  }
}
`, packageName, packageName, packageName, packageName)
}

// GenerateDependencyInjector generates dependency_injector.dart
func GenerateDependencyInjector() string {
	return `import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:logger/logger.dart';
import 'package:pine/di/dependency_injector_helper.dart';
import 'package:pine/utils/mapper.dart';
import 'package:pretty_dio_logger/pretty_dio_logger.dart';
import 'package:provider/provider.dart';
import 'package:provider/single_child_widget.dart';

part 'blocs.dart';
part 'mappers.dart';
part 'providers.dart';
part 'repositories.dart';

class DependencyInjector extends StatelessWidget {
  const DependencyInjector({super.key, required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) => DependencyInjectorHelper(
      repositories: _repositories,
      providers: _providers,
      blocs: _blocs,
      mappers: _mappers,
      child: child);
}
`
}

// GenerateBlocs generates blocs.dart
func GenerateBlocs() string {
	return `part of 'dependency_injector.dart';

final List<BlocProvider> _blocs = [
  // Add your BLoCs here
];
`
}

// GenerateMappers generates mappers.dart
func GenerateMappers() string {
	return `part of 'dependency_injector.dart';

final List<SingleChildWidget> _mappers = [
  // Add your mappers here
];
`
}

// GenerateProviders generates providers.dart
func GenerateProviders() string {
	return `part of 'dependency_injector.dart';

final List<SingleChildWidget> _providers = [
  Provider<Logger>(create: (_) => Logger()),

  Provider<PrettyDioLogger>(
      create: (_) => PrettyDioLogger(
          requestBody: true, compact: true, requestHeader: true)),

  Provider<Dio>(
      create: (context) => Dio()
        ..interceptors
            .addAll([if (kDebugMode) context.read<PrettyDioLogger>()])),

  Provider<FlutterSecureStorage>(
    create: (_) => const FlutterSecureStorage(),
  ),
];
`
}

// GenerateRepositories generates repositories.dart
func GenerateRepositories() string {
	return `part of 'dependency_injector.dart';

final List<RepositoryProvider> _repositories = [
  // Add your repositories here
];
`
}

// GenerateLightTheme generates light_theme.dart
func GenerateLightTheme() string {
	return `import 'package:flutter/material.dart';

class LightTheme {
  static ThemeData get make {
    return ThemeData(
      colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
      useMaterial3: true,
    );
  }
}
`
}

// GenerateAppRouter generates app_router.dart
func GenerateAppRouter(packageName string) string {
	return fmt.Sprintf(`import 'package:auto_route/auto_route.dart';
import 'package:%s/routers/app_router.gr.dart';

@AutoRouterConfig(
  replaceInRouteName: 'Page,Route',
)
class AppRouter extends RootStackRouter {
  @override
  List<AutoRoute> get routes => [
        AutoRoute(page: HomeRoute.page, initial: true),
      ];
}
`, packageName)
}

// GenerateL10n generates localization files
func GenerateL10n(locale string, appName string) string {
	return fmt.Sprintf(`{
    "appTitle": "%s",
    "@appTitle": {
        "description": "The title of the application"
    },
    "hello": "Hello",
    "@hello": {
        "description": "A greeting"
    }
}
`, appName)
}

// GenerateL10nYaml generates l10n.yaml configuration file
func GenerateL10nYaml() string {
	return `synthetic-package: false
arb-dir: lib/l10n
template-arb-file: app_en.arb
output-localization-file: app_localizations.dart
`
}
