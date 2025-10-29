package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileWriter handles file operations
type FileWriter struct {
	baseDir string
}

// NewFileWriter creates a new file writer
func NewFileWriter(baseDir string) *FileWriter {
	return &FileWriter{baseDir: baseDir}
}

// WriteFile writes content to a file, creating directories as needed
func (fw *FileWriter) WriteFile(relativePath string, content string) error {
	fullPath := filepath.Join(fw.baseDir, relativePath)

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	return nil
}

// EnsureDir creates a directory if it doesn't exist
func (fw *FileWriter) EnsureDir(relativePath string) error {
	fullPath := filepath.Join(fw.baseDir, relativePath)
	return os.MkdirAll(fullPath, 0755)
}

// PathExists checks if a path exists
func (fw *FileWriter) PathExists(relativePath string) bool {
	fullPath := filepath.Join(fw.baseDir, relativePath)
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

// ReadFile reads a file
func (fw *FileWriter) ReadFile(relativePath string) (string, error) {
	fullPath := filepath.Join(fw.baseDir, relativePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}
	return string(content), nil
}

// DeletePath deletes a file or directory
func (fw *FileWriter) DeletePath(relativePath string) error {
	fullPath := filepath.Join(fw.baseDir, relativePath)
	return os.RemoveAll(fullPath)
}

// GetFullPath returns the full path for a relative path
func (fw *FileWriter) GetFullPath(relativePath string) string {
	return filepath.Join(fw.baseDir, relativePath)
}
