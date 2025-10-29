package config

// ProjectConfig holds all configuration for project generation
type ProjectConfig struct {
	// Basic info
	ProjectName      string
	OrganizationName string
	Description      string
	TargetDirectory  string // Directory where to create the project

	// Backend options
	UseFirebase bool
	UseSupabase bool

	// Features
	EnableNotifications bool
	NotificationService string // "fcm" or "onesignal"

	// Models to generate
	Models []ModelConfig

	// Screens to generate
	GenerateLoginScreen    bool
	GenerateHomeScreen     bool
	GenerateProfileScreen  bool
	GenerateSettingsScreen bool

	// Additional options
	Force bool
}

// ModelConfig represents a model to generate
type ModelConfig struct {
	Name     string
	JSONData map[string]interface{}
	Endpoint string
}

// FirebaseConfig holds Firebase-specific configuration
type FirebaseConfig struct {
	EnableAuth         bool
	EnableFirestore    bool
	EnableStorage      bool
	EnableMessaging    bool
	EnableAnalytics    bool
	EnableCrashlytics  bool
}

// SupabaseConfig holds Supabase-specific configuration
type SupabaseConfig struct {
	ProjectURL string
	AnonKey    string
	EnableAuth bool
	EnableDB   bool
	EnableStorage bool
}

// DefaultProjectConfig returns a new config with sensible defaults
func DefaultProjectConfig() *ProjectConfig {
	return &ProjectConfig{
		OrganizationName:       "com.example",
		TargetDirectory:        ".",
		EnableNotifications:    false,
		GenerateLoginScreen:    true,
		GenerateHomeScreen:     true,
		GenerateProfileScreen:  false,
		GenerateSettingsScreen: false,
		Models:                 []ModelConfig{},
		Force:                  false,
	}
}
