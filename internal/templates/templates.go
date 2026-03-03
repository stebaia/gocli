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
		"logger: ^2.6.2",
		"flutter_secure_storage: ^10.0.0",
		"flutter_bloc: ^9.1.1",
		"hydrated_bloc: ^10.1.1",
		"equatable: ^2.0.8",
		"font_awesome_flutter: ^10.12.0",
		"pine: ^1.0.4",
		"provider: ^6.1.5+1",
		"retrofit: ^4.9.2",
		"dio: ^5.9.2",
		"pretty_dio_logger: ^1.4.0",
		"auto_route: ^11.1.0",
		"cached_network_image: ^3.4.1",
		"json_annotation: ^4.11.0",
		"sqlite3: ^3.0.0",
		"path_provider: ^2.1.5",
		"path: ^1.9.1",
		"shared_preferences: ^2.5.4",
		"intl: ^0.20.2",
		"google_fonts: ^8.0.2",
		"flutter_local_notifications: ^20.1.0",
	}

	// Add Firebase dependencies
	if cfg.UseFirebase {
		dependencies = append(dependencies,
			"firebase_core: ^4.5.0",
			"firebase_auth: ^6.2.0",
			"cloud_firestore: ^6.1.3",
			"google_sign_in: ^7.2.0",
		)

		if cfg.EnableNotifications {
			dependencies = append(dependencies, "firebase_messaging: ^16.1.2")
		}
	}

	// Add Supabase dependencies
	if cfg.UseSupabase {
		dependencies = append(dependencies, "supabase_flutter: ^2.12.0")
	}

	devDependencies := []string{
		"flutter_test:",
		"  sdk: flutter",
		"flutter_lints: ^6.0.0",
		"build_runner: ^2.11.1",
		"bloc_test: ^10.0.0",
		"retrofit_generator: ^10.2.3",
		"auto_route_generator: ^10.5.0",
		"http_mock_adapter: ^0.6.1",
		"data_fixture_dart: ^3.0.0",
		"mockito: ^5.6.3",
		"json_serializable: ^6.13.0",
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

// GenerateClaudeMd generates CLAUDE.md with project guidelines for LLMs
func GenerateClaudeMd() string {
	return `# Flutter Project Guidelines — fline architecture

This file contains the rules that **must be strictly followed** in every Flutter project generated with ` + "`fline`" + `. Every instruction here takes absolute priority over any general Flutter convention.

---

## Folder structure

` + "```" + `
lib/
├── di/                          # Dependency injection
│   ├── dependency_injector.dart
│   ├── blocs.dart
│   ├── mappers.dart
│   ├── providers.dart
│   └── repositories.dart
├── l10n/                        # Localizations
│   ├── app_en.arb
│   └── app_it.arb
├── mappers/                     # Mappers between DTO and domain model
├── model/                       # Domain models (Equatable + json_serializable)
├── network/
│   ├── interceptor/             # Dio interceptors
│   └── service/                 # Retrofit services
├── repositories/                # Repository pattern
├── routers/                     # auto_route
├── state_management/
│   ├── bloc/                    # BLoC (complex events, async streams)
│   ├── cubit/                   # Cubit (simple logic)
│   └── provider/                # Provider (global state without business logic)
├── theme/                       # App theme
│   └── light_theme.dart
├── ui/                          # Screens and widgets
│   ├── <screen_name>/
│   │   ├── <screen_name>_page.dart       # Screen entry point (annotated @RoutePage)
│   │   └── widgets/                      # Screen-specific widgets
│   │       └── <widget_name>.dart
├── utils/                       # Utilities and helpers
├── app.dart
└── main.dart
` + "```" + `

---

## Absolute rules — NO EXCEPTIONS

### 1. FORBIDDEN: setState

` + "`setState`" + ` is **strictly forbidden** across the entire codebase.
Any UI state must be managed through **BLoC**, **Cubit**, or **Provider**.

` + "```dart" + `
// ❌ FORBIDDEN
setState(() => _isLoading = true);

// ✅ CORRECT — use a Cubit
class LoginCubit extends Cubit<LoginState> {
  LoginCubit() : super(LoginInitial());

  Future<void> login(String email, String password) async {
    emit(LoginLoading());
    try {
      // authentication logic
      emit(LoginSuccess());
    } catch (e) {
      emit(LoginError(e.toString()));
    }
  }
}
` + "```" + `

### 2. FORBIDDEN: functions that return Widgets

Creating methods that return ` + "`Widget`" + ` is never allowed — neither inside a class nor as global functions. Every reusable piece of UI **must be a separate widget** in the appropriate ` + "`widgets/`" + ` folder.

` + "```dart" + `
// ❌ FORBIDDEN
Widget _buildHeader() {
  return Text('Header');
}

Widget buildButton(String label) {
  return ElevatedButton(...);
}

// ✅ CORRECT — separate widget in the widgets/ folder
// lib/ui/home/widgets/home_header.dart
class HomeHeader extends StatelessWidget {
  const HomeHeader({super.key});

  @override
  Widget build(BuildContext context) {
    return Text('Header');
  }
}
` + "```" + `

### 3. State management: when to use what

| Use case | Tool |
|---|---|
| Complex business logic, multiple events, stream transformations | **BLoC** |
| Simple logic (toggle, counter, form state) | **Cubit** |
| Global state without business logic (e.g. current user, theme) | **Provider** |
| Local UI state **without** business logic | ` + "`StatelessWidget`" + ` + props |

**Never use ` + "`StatefulWidget`" + ` to manage state that depends on business logic.**

---

## Internationalization (l10n)

Every user-visible string **must** be localized via ` + "`AppLocalizations`" + `.
The default supported languages are **Italian (it)** and **English (en)**.

### ARB files

- ` + "`lib/l10n/app_en.arb`" + ` — English strings (template)
- ` + "`lib/l10n/app_it.arb`" + ` — Italian strings

Every new string must be added to **both** files.

` + "```json" + `
// app_en.arb
{
  "welcomeMessage": "Welcome back",
  "@welcomeMessage": {
    "description": "Greeting on the login page"
  },
  "loginButton": "Sign In",
  "@loginButton": {
    "description": "Label of the login button"
  }
}

// app_it.arb
{
  "welcomeMessage": "Bentornato",
  "@welcomeMessage": {
    "description": "Saluto nella pagina di login"
  },
  "loginButton": "Accedi",
  "@loginButton": {
    "description": "Etichetta del pulsante di login"
  }
}
` + "```" + `

### Usage in code

` + "```dart" + `
// ❌ FORBIDDEN — hardcoded string
Text('Welcome back')
Text('Accedi')

// ✅ CORRECT
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

Text(AppLocalizations.of(context)!.welcomeMessage)
Text(AppLocalizations.of(context)!.loginButton)

// Recommended shorthand (add extension in utils/)
extension BuildContextL10n on BuildContext {
  AppLocalizations get l10n => AppLocalizations.of(this)!;
}

// Usage with extension
Text(context.l10n.welcomeMessage)
` + "```" + `

---

## Theming

The app theme is centralized in ` + "`lib/theme/`" + `. Hardcoded colors, fonts, or dimensions must never exist in the UI.

### Theme structure

` + "```dart" + `
// lib/theme/light_theme.dart
class LightTheme {
  static ThemeData get make {
    return ThemeData(
      colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
      useMaterial3: true,
    );
  }
}
` + "```" + `

### Theme rules

` + "```dart" + `
// ❌ FORBIDDEN — hardcoded values
Text('Title', style: TextStyle(fontSize: 24, color: Color(0xFF333333)))
Container(color: Colors.blue)

// ✅ CORRECT — use the theme
Text('Title', style: Theme.of(context).textTheme.headlineMedium)
Container(color: Theme.of(context).colorScheme.primary)
Container(color: Theme.of(context).colorScheme.surface)
` + "```" + `

To add custom colors or styles, extend the theme via ` + "`ThemeExtension`" + `.

---

## Layered architecture

The data flow strictly follows this order:

` + "```" + `
UI (Page/Widget)
    ↕  (BLoC/Cubit events & states)
State Management (BLoC / Cubit)
    ↕  (calls)
Repository
    ↕  (calls)
Network Service (Retrofit)
    ↕  (HTTP)
Backend API
` + "```" + `

### Rules per layer

**Model** (` + "`lib/model/`" + `)
- Must extend ` + "`Equatable`" + `
- Must use ` + "`@JsonSerializable()`" + ` with ` + "`json_annotation`" + `
- Immutable: all fields are ` + "`final`" + `
- No business logic inside

**Service** (` + "`lib/network/service/`" + `)
- Retrofit interface annotated with ` + "`@RestApi()`" + `
- Only HTTP endpoint declarations
- No logic

**Repository** (` + "`lib/repositories/`" + `)
- Depends on ` + "`Service`" + ` and ` + "`Logger`" + `
- Handles exceptions with ` + "`_logger.e(...)`" + ` and ` + "`rethrow`" + `
- May coordinate multiple services or local cache

**BLoC / Cubit** (` + "`lib/state_management/`" + `)
- Depends on Repository via constructor (dependency injection)
- Never accesses ` + "`Service`" + ` or ` + "`Dio`" + ` directly
- States must extend ` + "`Equatable`" + `

**Page** (` + "`lib/ui/<screen>/`" + `)
- Annotated with ` + "`@RoutePage()`" + `
- Contains no business logic
- Uses ` + "`BlocProvider`" + ` / ` + "`BlocBuilder`" + ` / ` + "`BlocListener`" + ` to access state
- Delegates UI to widgets in the ` + "`widgets/`" + ` subfolder

---

## Dependency Injection

DI is managed via the ` + "`pine`" + ` package through ` + "`DependencyInjector`" + `.
All providers, repositories, blocs, and mappers are registered in their respective ` + "`part of`" + ` files:

` + "```dart" + `
// lib/di/providers.dart
final List<SingleChildWidget> _providers = [
  Provider<Logger>(create: (_) => Logger()),
  Provider<Dio>(create: (context) => Dio()
    ..interceptors.add(context.read<PrettyDioLogger>())),
  Provider<MyService>(create: (context) => MyService(context.read<Dio>())),
];

// lib/di/repositories.dart
final List<RepositoryProvider> _repositories = [
  RepositoryProvider<MyRepository>(
    create: (context) => MyRepository(
      service: context.read<MyService>(),
      logger: context.read<Logger>(),
    ),
  ),
];

// lib/di/blocs.dart
final List<BlocProvider> _blocs = [
  BlocProvider<MyBloc>(
    create: (context) => MyBloc(
      repository: context.read<MyRepository>(),
    ),
  ),
];
` + "```" + `

---

## Routing with auto_route

All pages must be annotated with ` + "`@RoutePage()`" + ` and registered in ` + "`AppRouter`" + `.

` + "```dart" + `
// lib/routers/app_router.dart
@AutoRouterConfig(replaceInRouteName: 'Page,Route')
class AppRouter extends RootStackRouter {
  @override
  List<AutoRoute> get routes => [
    AutoRoute(page: HomeRoute.page, initial: true),
    AutoRoute(page: LoginRoute.page),
    AutoRoute(page: ProfileRoute.page),
  ];
}
` + "```" + `

Navigation inside Pages/Widgets:
` + "```dart" + `
// ✅ CORRECT — use AutoRoute
context.router.push(const ProfileRoute());
context.router.replace(const HomeRoute());
context.router.pop();

// ❌ FORBIDDEN — direct Navigator
Navigator.push(context, MaterialPageRoute(...));
Navigator.pushNamed(context, '/profile');
` + "```" + `

---

## Naming conventions

| Type | Convention | Example |
|---|---|---|
| File | ` + "`snake_case`" + ` | ` + "`user_profile_page.dart`" + ` |
| Class | ` + "`PascalCase`" + ` | ` + "`UserProfilePage`" + ` |
| Variable / method | ` + "`camelCase`" + ` | ` + "`fetchUserData()`" + ` |
| Constant | ` + "`camelCase`" + ` or ` + "`SCREAMING_SNAKE_CASE`" + ` | ` + "`defaultTimeout`" + ` |
| Screen | ` + "`Page`" + ` suffix | ` + "`LoginPage`" + `, ` + "`HomePage`" + ` |
| BLoC | ` + "`Bloc`" + ` suffix | ` + "`AuthBloc`" + ` |
| Cubit | ` + "`Cubit`" + ` suffix | ` + "`LoginCubit`" + ` |
| BLoC event | descriptive PascalCase | ` + "`FetchUsers`" + `, ` + "`DeleteUser`" + ` |
| BLoC state | state suffix | ` + "`UsersLoaded`" + `, ` + "`UserError`" + ` |
| Repository | ` + "`Repository`" + ` suffix | ` + "`UserRepository`" + ` |
| Service | ` + "`Service`" + ` suffix | ` + "`UserService`" + ` |
| Mapper | ` + "`Mapper`" + ` suffix | ` + "`UserMapper`" + ` |

---

## Main dependencies

| Package | Version | Purpose |
|---|---|---|
| ` + "`flutter_bloc`" + ` | ^9.1.1 | BLoC and Cubit |
| ` + "`hydrated_bloc`" + ` | ^10.1.1 | BLoC with persistence |
| ` + "`equatable`" + ` | ^2.0.8 | Object comparison |
| ` + "`provider`" + ` | ^6.1.5 | Provider pattern |
| ` + "`pine`" + ` | ^1.0.4 | DI helper |
| ` + "`auto_route`" + ` + ` + "`auto_route_generator`" + ` | ^11.1.0 / ^10.5.0 | Routing |
| ` + "`dio`" + ` | ^5.9.2 | HTTP client |
| ` + "`retrofit`" + ` + ` + "`retrofit_generator`" + ` | ^4.9.2 / ^10.2.3 | REST client codegen |
| ` + "`json_annotation`" + ` + ` + "`json_serializable`" + ` | ^4.11.0 / ^6.13.0 | JSON serialization |
| ` + "`flutter_secure_storage`" + ` | ^10.0.0 | Secure storage |
| ` + "`shared_preferences`" + ` | ^2.5.4 | User preferences |
| ` + "`cached_network_image`" + ` | ^3.4.1 | Images with cache |
| ` + "`logger`" + ` | ^2.6.2 | Logging |
| ` + "`google_fonts`" + ` | ^8.0.2 | Custom fonts |
| ` + "`intl`" + ` | ^0.20.2 | Internationalization |
| ` + "`build_runner`" + ` | ^2.11.1 | Code generation |

---

## Code generation

After adding or modifying models, services, or routers, run:

` + "```bash" + `
# Regenerate everything (json, retrofit, auto_route)
flutter pub run build_runner build --delete-conflicting-outputs

# Regenerate localizations
flutter gen-l10n
` + "```" + `

---

## Design reference from assets

If the ` + "`assets/images/`" + ` folder contains images (mockups, screenshots, UI designs), they **must** be used as the primary design reference for the entire project. This applies to both component structure and theming.

### Component extraction

Analyze every image in ` + "`assets/images/`" + ` and extract **all visible UI components** as separate widget classes. Do not approximate or skip elements — every button, card, input field, bottom sheet, list item, badge, avatar, or custom component visible in the designs must become its own widget file under ` + "`lib/ui/<screen>/widgets/`" + ` or a shared widget under ` + "`lib/ui/shared/widgets/`" + ` if reused across multiple screens.

` + "```" + `
// Example: if the design shows a custom card with an avatar, title and a tag badge,
// create three separate widgets:

lib/ui/shared/widgets/
├── user_avatar.dart        // the avatar component
├── tag_badge.dart          // the badge component
└── user_card.dart          // the card that composes the above two
` + "```" + `

Never inline a component visible in the design directly inside a Page's ` + "`build`" + ` method — always extract it.

### Theming from design

Colors, typography, border radii, spacing, and any visual style visible in the design images **must be reflected in the theme**. Do not hardcode values extracted from the design — register them in ` + "`lib/theme/`" + ` instead.

**Colors:** extract the primary, secondary, background, surface, and accent colors from the designs and define them in the ` + "`ColorScheme`" + `.

**Typography:** if the design uses specific font weights, sizes, or a particular font family, configure them in ` + "`TextTheme`" + ` using ` + "`google_fonts`" + ` if needed.

**Shape / border radius:** if the design uses rounded corners consistently, define a ` + "`ShapeBorder`" + ` or use ` + "`ThemeData.cardTheme`" + `, ` + "`ThemeData.inputDecorationTheme`" + `, etc.

` + "```dart" + `
// lib/theme/light_theme.dart
class LightTheme {
  static ThemeData get make {
    return ThemeData(
      colorScheme: ColorScheme.fromSeed(
        seedColor: const Color(0xFF4F46E5), // extracted from design
        primary: const Color(0xFF4F46E5),
        secondary: const Color(0xFF10B981),
        surface: const Color(0xFFF9FAFB),
      ),
      useMaterial3: true,
      textTheme: GoogleFonts.interTextTheme(), // if Inter is used in the design
      cardTheme: const CardTheme(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(16)), // from design
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12), // from design
        ),
      ),
    );
  }
}
` + "```" + `

### Priority rule

> If a design image is present in ` + "`assets/images/`" + `, it overrides any default or placeholder implementation. The generated UI **must match** the design as closely as possible. Generic scaffolding (e.g. the default ` + "`HomePage`" + ` or ` + "`LoginPage`" + `) must be replaced with components derived from the actual design.

---

## Pre-commit checklist

- [ ] No ` + "`setState`" + ` in the code
- [ ] No functions returning ` + "`Widget`" + ` (use dedicated widget classes)
- [ ] Every user-visible string is localized in both ` + "`app_en.arb`" + ` and ` + "`app_it.arb`" + `
- [ ] No hardcoded colors/fonts/dimensions (use the theme)
- [ ] Navigation via ` + "`auto_route`" + ` (no direct ` + "`Navigator`" + `)
- [ ] New BLoC/Cubit/Repositories registered in the DI
- [ ] ` + "`build_runner`" + ` run after changes to models/services/routers
- [ ] ` + "`flutter gen-l10n`" + ` run after changes to ARB files
`
}

// GenerateL10nYaml generates l10n.yaml configuration file
func GenerateL10nYaml() string {
	return `synthetic-package: false
arb-dir: lib/l10n
template-arb-file: app_en.arb
output-localization-file: app_localizations.dart
`
}
