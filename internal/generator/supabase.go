package generator

import (
	"fline-cli/internal/config"
	"fline-cli/internal/ui"
	"fline-cli/internal/utils"
)

// SupabaseGenerator handles Supabase integration
type SupabaseGenerator struct {
	config  *config.ProjectConfig
	writer  *utils.FileWriter
	flutter *utils.FlutterCLI
	logger  *ui.Logger
}

// NewSupabaseGenerator creates a new Supabase generator
func NewSupabaseGenerator(cfg *config.ProjectConfig, writer *utils.FileWriter, flutter *utils.FlutterCLI) *SupabaseGenerator {
	return &SupabaseGenerator{
		config:  cfg,
		writer:  writer,
		flutter: flutter,
		logger:  ui.NewLogger("supabase"),
	}
}

// Generate sets up Supabase
func (g *SupabaseGenerator) Generate() error {
	g.logger.Info("Setting up Supabase integration...")

	// Generate Supabase client
	if err := g.generateSupabaseClient(); err != nil {
		return err
	}

	// Generate auth service
	if err := g.generateAuthService(); err != nil {
		return err
	}

	g.logger.Success("Supabase integration configured")
	g.logger.Info("Don't forget to:")
	g.logger.Info("1. Create a Supabase project at https://supabase.com")
	g.logger.Info("2. Get your project URL and anon key")
	g.logger.Info("3. Update lib/utils/supabase_client.dart with your credentials")

	return nil
}

func (g *SupabaseGenerator) generateSupabaseClient() error {
	content := `import 'package:supabase_flutter/supabase_flutter.dart';

class SupabaseConfig {
  static const String supabaseUrl = 'YOUR_SUPABASE_URL';
  static const String supabaseAnonKey = 'YOUR_SUPABASE_ANON_KEY';

  static Future<void> initialize() async {
    await Supabase.initialize(
      url: supabaseUrl,
      anonKey: supabaseAnonKey,
    );
  }

  static SupabaseClient get client => Supabase.instance.client;
}
`

	return g.writer.WriteFile("lib/utils/supabase_client.dart", content)
}

func (g *SupabaseGenerator) generateAuthService() error {
	content := `import 'package:logger/logger.dart';
import 'package:supabase_flutter/supabase_flutter.dart';

class SupabaseAuthService {
  final SupabaseClient _client;
  final Logger _logger;

  SupabaseAuthService({
    required SupabaseClient client,
    required Logger logger,
  })  : _client = client,
        _logger = logger;

  // Get current user
  User? get currentUser => _client.auth.currentUser;

  // Auth state changes
  Stream<AuthState> get authStateChanges => _client.auth.onAuthStateChange;

  // Sign in with email and password
  Future<AuthResponse> signInWithEmailAndPassword({
    required String email,
    required String password,
  }) async {
    try {
      return await _client.auth.signInWithPassword(
        email: email,
        password: password,
      );
    } catch (e) {
      _logger.e('Sign in error', error: e);
      rethrow;
    }
  }

  // Register with email and password
  Future<AuthResponse> signUp({
    required String email,
    required String password,
  }) async {
    try {
      return await _client.auth.signUp(
        email: email,
        password: password,
      );
    } catch (e) {
      _logger.e('Sign up error', error: e);
      rethrow;
    }
  }

  // Sign out
  Future<void> signOut() async {
    try {
      await _client.auth.signOut();
    } catch (e) {
      _logger.e('Sign out error', error: e);
      rethrow;
    }
  }

  // Reset password
  Future<void> resetPassword(String email) async {
    try {
      await _client.auth.resetPasswordForEmail(email);
    } catch (e) {
      _logger.e('Password reset error', error: e);
      rethrow;
    }
  }
}
`

	if err := g.writer.WriteFile("lib/network/service/supabase_auth_service.dart", content); err != nil {
		return err
	}

	// Add note about providers
	providerContent := `
  // Supabase
  Provider<SupabaseClient>(create: (_) => Supabase.instance.client),
  Provider<SupabaseAuthService>(
    create: (context) => SupabaseAuthService(
      client: context.read<SupabaseClient>(),
      logger: context.read<Logger>(),
    ),
  ),
`
	g.logger.Warning("Remember to add Supabase providers to lib/di/providers.dart:")
	g.logger.Info(providerContent)

	return nil
}
