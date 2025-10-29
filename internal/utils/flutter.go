package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// FlutterCLI wraps Flutter CLI commands
type FlutterCLI struct {
	workingDir string
}

// NewFlutterCLI creates a new Flutter CLI wrapper
func NewFlutterCLI(workingDir string) *FlutterCLI {
	return &FlutterCLI{workingDir: workingDir}
}

// Create creates a new Flutter project
func (f *FlutterCLI) Create(projectName, org string, force bool) error {
	args := []string{
		"create",
		"--org", org,
		"--project-name", projectName,
	}

	if force {
		args = append(args, "--overwrite")
	}

	args = append(args, projectName)

	cmd := exec.Command("flutter", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("flutter create failed: %s\n%s", err, string(output))
	}

	return nil
}

// PubGet runs flutter pub get
func (f *FlutterCLI) PubGet() error {
	return f.runCommand("pub", "get")
}

// PubAdd adds a package
func (f *FlutterCLI) PubAdd(packages ...string) error {
	args := append([]string{"pub", "add"}, packages...)
	return f.runCommand(args...)
}

// PubAddDev adds a dev package
func (f *FlutterCLI) PubAddDev(packages ...string) error {
	args := append([]string{"pub", "add", "--dev"}, packages...)
	return f.runCommand(args...)
}

// GenL10n generates localizations
func (f *FlutterCLI) GenL10n() error {
	return f.runCommand("gen-l10n")
}

// BuildRunnerBuild runs build_runner
func (f *FlutterCLI) BuildRunnerBuild() error {
	return f.runCommand("pub", "run", "build_runner", "build", "--delete-conflicting-outputs")
}

// DartRun runs a dart command
func (f *FlutterCLI) DartRun(args ...string) error {
	fullArgs := append([]string{"pub", "run"}, args...)
	return f.runCommand(fullArgs...)
}

// runCommand executes a flutter command
func (f *FlutterCLI) runCommand(args ...string) error {
	cmd := exec.Command("flutter", args...)
	cmd.Dir = f.workingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flutter %s failed: %s\n%s",
			strings.Join(args, " "), err, string(output))
	}

	return nil
}

// CheckFlutterInstalled checks if Flutter is installed
func CheckFlutterInstalled() error {
	cmd := exec.Command("flutter", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("flutter is not installed or not in PATH")
	}
	return nil
}
