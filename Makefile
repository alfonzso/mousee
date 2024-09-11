# Name of the Go binary
BINARY_NAME=mousee.exe

# Go source file
SOURCE_FILE=main.go

# Default target: remove, build, and run
all: kill remove build run

# Remove the binary if it exists
remove:
		@if [ -f $(BINARY_NAME) ]; then \
				echo "Removing existing binary $(BINARY_NAME)..."; \
				rm $(BINARY_NAME); \
		else \
				echo "No existing binary found."; \
		fi

# Kill: Kill the process if it's running
kill:
	@echo "Killing $(BINARY_NAME) process..."
	@taskkill //F //IM $(BINARY_NAME) || echo 1

# Build the Go binary
build:
		@echo "Building Go binary..."
		go build .
		# go build -o $(BINARY_NAME) $(SOURCE_FILE)

# Run the Go binary
run:
		@echo "Running $(BINARY_NAME)..."
		./$(BINARY_NAME)

# Clean up the binary
clean:
		@if [ -f $(BINARY_NAME) ]; then \
				echo "Cleaning up $(BINARY_NAME)..."; \
				rm $(BINARY_NAME); \
		else \
				echo "No binary to clean."; \
		fi
