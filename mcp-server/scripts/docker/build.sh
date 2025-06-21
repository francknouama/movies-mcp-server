#!/bin/bash

# ==============================================================================
# Docker Build Script for Movies MCP Server
# ==============================================================================

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
IMAGE_NAME="movies-mcp-server"
REGISTRY="${REGISTRY:-ghcr.io/your-org}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Help function
show_help() {
    cat << EOF
Docker Build Script for Movies MCP Server

Usage: $0 [OPTIONS]

OPTIONS:
    -t, --tag TAG           Tag for the Docker image (default: latest)
    -p, --push              Push image to registry after build
    -s, --security          Use security-hardened Dockerfile
    -c, --cache             Use Docker buildx cache
    --no-cache              Build without cache
    --platform PLATFORMS   Target platforms (default: linux/amd64)
    -h, --help              Show this help message

EXAMPLES:
    $0                      # Build with default settings
    $0 -t v1.0.0 -p         # Build version 1.0.0 and push
    $0 -s --no-cache        # Security build without cache
    $0 --platform linux/amd64,linux/arm64  # Multi-platform build

ENVIRONMENT VARIABLES:
    REGISTRY               Docker registry (default: ghcr.io/your-org)
    DOCKER_BUILDKIT        Enable BuildKit (default: 1)
EOF
}

# Default values
TAG="latest"
PUSH=false
SECURITY=false
USE_CACHE=true
NO_CACHE=false
PLATFORM="linux/amd64"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--tag)
            TAG="$2"
            shift 2
            ;;
        -p|--push)
            PUSH=true
            shift
            ;;
        -s|--security)
            SECURITY=true
            shift
            ;;
        -c|--cache)
            USE_CACHE=true
            shift
            ;;
        --no-cache)
            NO_CACHE=true
            USE_CACHE=false
            shift
            ;;
        --platform)
            PLATFORM="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validation
if ! command -v docker &> /dev/null; then
    log_error "Docker is not installed or not in PATH"
    exit 1
fi

if ! command -v git &> /dev/null; then
    log_warning "Git is not installed. Version info will be limited."
fi

# Change to project root
cd "$PROJECT_ROOT"

# Get build metadata
VERSION="${TAG}"
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
if command -v git &> /dev/null && git rev-parse --git-dir > /dev/null 2>&1; then
    GIT_COMMIT=$(git rev-parse --short HEAD)
    GIT_DIRTY=$(git diff --quiet || echo "-dirty")
    GIT_COMMIT="${GIT_COMMIT}${GIT_DIRTY}"
else
    GIT_COMMIT="unknown"
fi

# Determine Dockerfile
DOCKERFILE="Dockerfile"
if [[ "$SECURITY" == "true" ]]; then
    DOCKERFILE="Dockerfile.production"
    log_info "Using security-hardened Dockerfile"
fi

# Full image name
FULL_IMAGE_NAME="${REGISTRY}/${IMAGE_NAME}:${TAG}"

# Build arguments
BUILD_ARGS=(
    "--build-arg" "VERSION=${VERSION}"
    "--build-arg" "BUILD_TIME=${BUILD_TIME}"
    "--build-arg" "GIT_COMMIT=${GIT_COMMIT}"
)

# Cache arguments
if [[ "$USE_CACHE" == "true" && "$NO_CACHE" == "false" ]]; then
    BUILD_ARGS+=(
        "--cache-from" "${REGISTRY}/${IMAGE_NAME}:buildcache"
        "--cache-to" "type=registry,ref=${REGISTRY}/${IMAGE_NAME}:buildcache,mode=max"
    )
elif [[ "$NO_CACHE" == "true" ]]; then
    BUILD_ARGS+=("--no-cache")
fi

# Platform arguments
BUILD_ARGS+=("--platform" "$PLATFORM")

# Enable BuildKit
export DOCKER_BUILDKIT=1

log_info "Building Docker image..."
log_info "Image: ${FULL_IMAGE_NAME}"
log_info "Version: ${VERSION}"
log_info "Build Time: ${BUILD_TIME}"
log_info "Git Commit: ${GIT_COMMIT}"
log_info "Platform: ${PLATFORM}"
log_info "Dockerfile: ${DOCKERFILE}"

# Build the image
if docker buildx build \
    "${BUILD_ARGS[@]}" \
    -f "$DOCKERFILE" \
    -t "$FULL_IMAGE_NAME" \
    -t "${REGISTRY}/${IMAGE_NAME}:latest" \
    .; then
    
    log_success "Docker image built successfully"
    
    # Show image info
    log_info "Image size: $(docker images --format "table {{.Size}}" "$FULL_IMAGE_NAME" | tail -n 1)"
    
    # Security scan (if available)
    if command -v trivy &> /dev/null; then
        log_info "Running security scan..."
        trivy image --exit-code 0 --severity HIGH,CRITICAL "$FULL_IMAGE_NAME"
    fi
    
    # Push if requested
    if [[ "$PUSH" == "true" ]]; then
        log_info "Pushing image to registry..."
        if docker push "$FULL_IMAGE_NAME" && docker push "${REGISTRY}/${IMAGE_NAME}:latest"; then
            log_success "Image pushed successfully"
        else
            log_error "Failed to push image"
            exit 1
        fi
    fi
    
    log_success "Build completed successfully!"
    log_info "Image: ${FULL_IMAGE_NAME}"
    
    # Usage instructions
    echo
    log_info "To run the container:"
    echo "  docker run --rm -it ${FULL_IMAGE_NAME}"
    echo
    log_info "To run with docker-compose:"
    echo "  docker-compose up"
    
else
    log_error "Docker build failed"
    exit 1
fi