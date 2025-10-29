# 🚀 Fline CLI

<div align="center">

**A powerful Flutter project generator with Pine architecture**

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Flutter](https://img.shields.io/badge/Flutter-3.0+-02569B?style=flat&logo=flutter)](https://flutter.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

</div>

## ✨ Features

Fline CLI helps you create professional Flutter projects with:

- 🏗️ **Clean Architecture** - Service, Repository, BLoC, UI layers
- 🔥 **Firebase Integration** - Auth, Firestore, Storage, Messaging
- 💧 **Supabase Integration** - Auth, Database, Storage
- 🔔 **Push Notifications** - FCM setup ready
- 📦 **Model Generation** - From JSON to complete features
- 🎨 **Example Screens** - Login, Home, Profile, Settings
- 🌈 **Beautiful CLI** - Interactive wizard with colors
- ⚡ **Fast & Modern** - Built with Go for speed

## 📋 Prerequisites

- [Flutter](https://flutter.dev/docs/get-started/install) installed
- [Go](https://golang.org/dl/) 1.24+ (for building from source)

## 🔧 Installation

### From Source

```bash
git clone https://github.com/yourusername/fline-cli
cd fline-cli
make install
```

### Usage without installation

```bash
go run main.go [command]
```

## 🎯 Commands

### `fline create` - Create a New Project

Create a new Flutter project with interactive wizard:

```bash
fline create
```

Or use flags for non-interactive mode:

```bash
fline create --name my_app --org com.example --firebase
```

**What you get:**
- ✅ Flutter project with Pine architecture
- ✅ Dependency injection setup
- ✅ Router configuration (auto_route)
- ✅ Localization (l10n) ready
- ✅ Theme configuration
- ✅ Optional Firebase/Supabase
- ✅ Optional example screens
- ✅ Optional models from JSON

### `fline generate` - Generate Features

Generate service, repository, and BLoC for a feature:

```bash
# Interactive mode
fline generate

# With arguments
fline generate user
fline generate product --type service
```

**Options:**
- `--name, -n`: Feature name
- `--type, -t`: Component type (all, service, repository, bloc)

### `fline model` - Generate from JSON

Create complete feature from JSON data:

```bash
# Interactive mode
fline model

# From JSON string
fline model user --json '{"id":1,"name":"John","email":"john@example.com"}' --endpoint /api/users

# From JSON file
fline model product --json-file product.json
```

**Generates:**
- 📄 Model with json_annotation
- 🌐 Retrofit service
- 💾 Repository with error handling
- 🎯 BLoC with CRUD operations (Create, Read, Update, Delete)

## 📁 Project Structure

```
your_project/
├── lib/
│   ├── di/                        # Dependency Injection
│   │   ├── dependency_injector.dart
│   │   ├── blocs.dart
│   │   ├── providers.dart
│   │   └── repositories.dart
│   ├── l10n/                      # Localizations
│   ├── model/                     # Data models
│   ├── network/
│   │   └── service/               # API services (Retrofit)
│   ├── repositories/              # Data repositories
│   ├── routers/                   # Navigation (auto_route)
│   ├── state_management/
│   │   └── bloc/                  # BLoC state management
│   ├── theme/                     # App theming
│   ├── ui/                        # UI screens & widgets
│   ├── utils/                     # Utilities
│   ├── app.dart
│   └── main.dart
├── pubspec.yaml
└── ...
```

## 🔥 Firebase Setup

After generating a project with Firebase:

1. Install FlutterFire CLI:
```bash
dart pub global activate flutterfire_cli
```

2. Configure Firebase:
```bash
cd your_project
flutterfire configure
```

3. Follow the prompts to select/create a Firebase project

## 💧 Supabase Setup

After generating a project with Supabase:

1. Create a project at [supabase.com](https://supabase.com)

2. Get your Project URL and anon key

3. Update `lib/utils/supabase_client.dart`:
```dart
static const String supabaseUrl = 'YOUR_PROJECT_URL';
static const String supabaseAnonKey = 'YOUR_ANON_KEY';
```

## 🎨 Example Screens

Fline CLI can generate these example screens:

- **Login Screen** - Email/password authentication UI
- **Home Screen** - Dashboard with quick actions
- **Profile Screen** - User profile management
- **Settings Screen** - App settings and preferences

All screens are:
- ✅ Material Design 3 compliant
- ✅ Responsive
- ✅ Well-structured
- ✅ Ready to customize

## 🛠️ Development Workflow

After creating a project:

```bash
cd your_project

# Install dependencies
flutter pub get

# Generate code (routes, models, etc.)
flutter pub run build_runner build --delete-conflicting-outputs

# Generate localizations
flutter gen-l10n

# Run the app
flutter run
```

## 📖 Examples

### Create a complete project

```bash
fline create
# Follow the interactive wizard:
# 1. Enter project name: my_awesome_app
# 2. Enter organization: com.mycompany
# 3. Select backend: Firebase
# 4. Enable notifications: Yes
# 5. Add models: Yes, paste JSON
# 6. Select screens: Login, Home, Profile
```

### Generate a user feature

```bash
cd my_awesome_app
fline model user --json '{"id":1,"name":"John Doe","email":"john@example.com","avatar":"url"}' --endpoint /api/users
```

This generates:
- `lib/model/user.dart`
- `lib/network/service/user_service.dart`
- `lib/repositories/user_repository.dart`
- `lib/state_management/bloc/user/user_bloc.dart`
- `lib/state_management/bloc/user/user_event.dart`
- `lib/state_management/bloc/user/user_state.dart`

### Generate just a service

```bash
fline generate post --type service
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Pine Package](https://pub.dev/packages/pine) - Flutter architecture helpers
- [Charm.sh](https://charm.sh/) - Beautiful CLI tools
- [Cobra](https://cobra.dev/) - CLI framework

## 💬 Support

- 🐛 Report bugs via [GitHub Issues](https://github.com/yourusername/fline-cli/issues)
- 💡 Feature requests welcome
- ⭐ Star the project if you like it!

---

<div align="center">
Made with ❤️ for Flutter developers
</div>
