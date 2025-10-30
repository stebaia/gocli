package utils

import (
	"github.com/iancoleman/strcase"
	"strings"
)

// NamingHelper provides case conversion utilities
type NamingHelper struct {
	original string
}

// NewNamingHelper creates a new naming helper
func NewNamingHelper(name string) *NamingHelper {
	return &NamingHelper{original: name}
}

// SnakeCase returns snake_case version
func (n *NamingHelper) SnakeCase() string {
	return strcase.ToSnake(n.original)
}

// PascalCase returns PascalCase version
func (n *NamingHelper) PascalCase() string {
	return strcase.ToCamel(n.original)
}

// CamelCase returns camelCase version
func (n *NamingHelper) CamelCase() string {
	return strcase.ToLowerCamel(n.original)
}

// KebabCase returns kebab-case version
func (n *NamingHelper) KebabCase() string {
	return strcase.ToKebab(n.original)
}

// ScreamingSnakeCase returns SCREAMING_SNAKE_CASE version
func (n *NamingHelper) ScreamingSnakeCase() string {
	return strcase.ToScreamingSnake(n.original)
}

// Original returns the original string
func (n *NamingHelper) Original() string {
	return n.original
}

// Reserved Dart/Flutter package names that cannot be used as project names
var reservedNames = []string{
	"test", "flutter", "dart", "build", "integration_test",
	"flutter_test", "flutter_driver", "sky_engine",
}

// ValidateProjectName checks if a project name is valid for Flutter
func ValidateProjectName(name string) error {
	if name == "" {
		return &ValidationError{Field: "name", Message: "project name cannot be empty"}
	}

	if strings.Contains(name, " ") {
		return &ValidationError{Field: "name", Message: "project name cannot contain spaces"}
	}

	if strings.Contains(name, "-") {
		return &ValidationError{Field: "name", Message: "project name should use underscores instead of dashes"}
	}

	// Check if starts with letter
	if !isLetter(rune(name[0])) {
		return &ValidationError{Field: "name", Message: "project name must start with a letter"}
	}

	// Check if name is reserved
	for _, reserved := range reservedNames {
		if strings.EqualFold(name, reserved) {
			return &ValidationError{
				Field:   "name",
				Message: "'" + name + "' is a reserved name and cannot be used. Please choose a different name",
			}
		}
	}

	return nil
}

// ValidateOrganization checks if an organization name is valid
func ValidateOrganization(org string) error {
	if org == "" {
		return &ValidationError{Field: "organization", Message: "organization cannot be empty"}
	}

	parts := strings.Split(org, ".")
	if len(parts) < 2 {
		return &ValidationError{Field: "organization", Message: "organization must be in format: com.example"}
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
