
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

clean:
	rm -f $(TARGET)

up:
	go build -o ./bin/ ./cmd/vcli...
