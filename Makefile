# Define the output directory
DIST_DIR := dist

# Define Go options
GO_BUILD_FLAGS := CGO_ENABLED=0

# Define the target platforms
PLATFORMS = \
	windows/386 \
    windows/amd64 \
	windows/arm \
    windows/arm64 \
    linux/amd64 \
	linux/arm \
    linux/arm64 \
    linux/386 \
    darwin/amd64 \
    darwin/arm64

default: test clean build zip

# Define the build process
build: $(PLATFORMS)
	@echo "Build complete!"

# Cross-compile for each platform
$(PLATFORMS):
	@echo "Building for $@..."
	GOOS=$(word 1, $(subst /, ,$@)) GOARCH=$(word 2, $(subst /, ,$@)) $(GO_BUILD_FLAGS) go build -o $(DIST_DIR)/otpmini_$(word 1, $(subst /, ,$@))_$(word 2, $(subst /, ,$@))$(if $(filter windows,$(word 1, $(subst /, ,$@))),.exe)

# Clean the dist directory
clean:
	@echo "Cleaning dist directory..."
	rm -rf $(DIST_DIR)

# Clean the dist directory
test:
	@echo "Testing..."
	go test . -v

# Create the dist directory if it doesn't exist
$(DIST_DIR):
	mkdir -p $(DIST_DIR)

# Zip the output files
zip: $(DIST_DIR)
	@echo "Creating zip files..."
	@for file in $(shell ls $(DIST_DIR)); do \
		zip -j $(DIST_DIR)/$$file.zip $(DIST_DIR)/$$file; \
	done

.PHONY: build clean test zip
