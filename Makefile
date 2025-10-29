.PHONY: build install uninstall clean help

# Binary name
BINARY_NAME=fline

# Installation directory
INSTALL_PATH=/usr/local/bin

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) main.go
	@echo "✓ Build complete: ./$(BINARY_NAME)"

# Install the binary globally
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo mv $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ $(BINARY_NAME) installed successfully!"
	@echo "You can now run '$(BINARY_NAME)' from anywhere"

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ $(BINARY_NAME) uninstalled"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@echo "✓ Clean complete"

# Show help
help:
	@echo "Fline CLI - Makefile commands:"
	@echo ""
	@echo "  make build      - Build the binary locally"
	@echo "  make install    - Build and install globally (requires sudo)"
	@echo "  make uninstall  - Remove the installed binary"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make help       - Show this help message"
