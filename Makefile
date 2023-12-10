# Makefile for your Go project

# Go environment variables
export GOOS = linux
export GOARCH = amd64

# Go compiler and build flags
GO = go
GOFLAGS = -tags netgo -installsuffix netgo -ldflags="-w -s"

# Source files
SRC = $(wildcard *.go)

# Target executable
TARGET = bin/vcli


# Default target
all: $(TARGET)

# Rule to build the target executable
$(TARGET): $(SRC)
	$(GO) build $(GOFLAGS) -o $@

# Clean rule
clean:
	rm -f $(TARGET)

up:
	go build -o ./bin/ ./...

publish:
	${env:GOOS}='linux';${env:GOARCH}='amd64'; & 'go build -tags netgo -installsuffix netgo -ldflags="-w -s" -o bin/ .\...' 