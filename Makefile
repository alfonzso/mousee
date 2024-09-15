# Name of the Go binary
BINARY_NAME=mousee.exe

# Go source file
SOURCE_FILE=main.go

# Default target: remove, build, and run
all: remove build run

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
# EPOCH=(powershell.exe -c 'Get-Date -Date ((Get-Date).DateTime) -UFormat %s'); \
# @echo $(shell -c 'Get-Date -Date ((Get-Date).DateTime) -UFormat %s')
# EPOCH=$$(powershell.exe -c 'Get-Date -Date ((Get-Date).DateTime) -UFormat %s') ; \

build:
	SHA=$$(powershell.exe -c 'git rev-parse --short HEAD') ; \
	MMDDHHMM=$$(powershell.exe -c 'Get-Date -Format MMdd.hhmm') ; \
	go build -o $(BINARY_NAME) -ldflags "-X common.AppVersion=v$$SHA.$$MMDDHHMM -X common.AppName=$(BINARY_NAME)" .

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
