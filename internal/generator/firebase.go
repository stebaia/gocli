package generator

import (
	"fmt"
	"fline-cli/internal/config"
	"fline-cli/internal/ui"
	"fline-cli/internal/utils"
)

// FirebaseGenerator handles Firebase integration
type FirebaseGenerator struct {
	config  *config.ProjectConfig
	writer  *utils.FileWriter
	flutter *utils.FlutterCLI
	logger  *ui.Logger
}

// NewFirebaseGenerator creates a new Firebase generator
func NewFirebaseGenerator(cfg *config.ProjectConfig, writer *utils.FileWriter, flutter *utils.FlutterCLI) *FirebaseGenerator {
	return &FirebaseGenerator{
		config:  cfg,
		writer:  writer,
		flutter: flutter,
		logger:  ui.NewLogger("firebase"),
	}
}

// Generate sets up Firebase
func (g *FirebaseGenerator) Generate() error {
	g.logger.Info("Setting up Firebase integration...")

	// Add Firebase packages
	packages := []string{
		"firebase_core",
		"firebase_auth",
		"cloud_firestore",
	}

	if g.config.EnableNotifications {
		packages = append(packages, "firebase_messaging")
	}

	// Note: packages are already in pubspec, so we don't need to add them again

	// Generate Firebase initialization
	if err := g.generateFirebaseInit(); err != nil {
		return err
	}

	// Generate auth service if needed
	if err := g.generateAuthService(); err != nil {
		return err
	}

	g.logger.Success("Firebase integration configured")
	g.logger.Info("Don't forget to:")
	g.logger.Info("1. Create a Firebase project at https://console.firebase.google.com")
	g.logger.Info("2. Run: flutterfire configure")
	g.logger.Info("3. Follow the setup instructions")

	return nil
}

func (g *FirebaseGenerator) generateFirebaseInit() error {
	content := `import 'package:firebase_core/firebase_core.dart';
import 'firebase_options.dart';

class FirebaseInitializer {
  static Future<void> initialize() async {
    await Firebase.initializeApp(
      options: DefaultFirebaseOptions.currentPlatform,
    );
  }
}
`

	return g.writer.WriteFile("lib/utils/firebase_initializer.dart", content)
}

func (g *FirebaseGenerator) generateAuthService() error {
	content := fmt.Sprintf(`import 'package:firebase_auth/firebase_auth.dart';
import 'package:logger/logger.dart';

class AuthService {
  final FirebaseAuth _auth;
  final Logger _logger;

  AuthService({
    required FirebaseAuth auth,
    required Logger logger,
  })  : _auth = auth,
        _logger = logger;

  // Get current user
  User? get currentUser => _auth.currentUser;

  // Auth state changes
  Stream<User?> get authStateChanges => _auth.authStateChanges();

  // Sign in with email and password
  Future<UserCredential> signInWithEmailAndPassword({
    required String email,
    required String password,
  }) async {
    try {
      return await _auth.signInWithEmailAndPassword(
        email: email,
        password: password,
      );
    } catch (e) {
      _logger.e('Sign in error', error: e);
      rethrow;
    }
  }

  // Register with email and password
  Future<UserCredential> registerWithEmailAndPassword({
    required String email,
    required String password,
  }) async {
    try {
      return await _auth.createUserWithEmailAndPassword(
        email: email,
        password: password,
      );
    } catch (e) {
      _logger.e('Registration error', error: e);
      rethrow;
    }
  }

  // Sign out
  Future<void> signOut() async {
    try {
      await _auth.signOut();
    } catch (e) {
      _logger.e('Sign out error', error: e);
      rethrow;
    }
  }

  // Reset password
  Future<void> resetPassword(String email) async {
    try {
      await _auth.sendPasswordResetEmail(email: email);
    } catch (e) {
      _logger.e('Password reset error', error: e);
      rethrow;
    }
  }
}
`)

	if err := g.writer.WriteFile("lib/network/service/auth_service.dart", content); err != nil {
		return err
	}

	// Add to providers
	return g.addAuthProvider()
}

func (g *FirebaseGenerator) addAuthProvider() error {
	// This is a simplified version - in production you'd want to properly insert into the existing file
	content := `
  // Firebase Auth
  Provider<FirebaseAuth>(create: (_) => FirebaseAuth.instance),
  Provider<AuthService>(
    create: (context) => AuthService(
      auth: context.read<FirebaseAuth>(),
      logger: context.read<Logger>(),
    ),
  ),
`
	g.logger.Warning("Remember to add Firebase providers to lib/di/providers.dart:")
	g.logger.Info(content)

	return nil
}
