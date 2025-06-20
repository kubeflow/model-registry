#!/bin/bash

# Script to run Model Registry E2E tests
# This script sets up the environment and runs the end-to-end tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

print_status "Running Model Registry E2E tests from: $PROJECT_DIR"

# Check if we're in a virtual environment
if [[ "$VIRTUAL_ENV" == "" ]]; then
    print_warning "Not running in a virtual environment. Consider activating one."
fi

# Check if required environment variables are set
if [[ -z "$MODEL_REGISTRY_HOST" ]]; then
    print_error "MODEL_REGISTRY_HOST environment variable is not set"
    print_status "Please set the required environment variables:"
    echo "  export MODEL_REGISTRY_HOST=your-model-registry-server.com"
    echo "  export MODEL_REGISTRY_PORT=8080  # optional, defaults to 8080"
    echo "  export MODEL_REGISTRY_SECURE=false  # optional, defaults to false"
    echo "  export MODEL_REGISTRY_TOKEN=your-auth-token"
    echo ""
    print_status "Or create an e2e_config.env file and source it:"
    echo "  cp tests/e2e_config.env.example tests/e2e_config.env"
    echo "  # Edit tests/e2e_config.env with your values"
    echo "  source tests/e2e_config.env"
    exit 1
fi

if [[ -z "$MODEL_REGISTRY_TOKEN" ]]; then
    print_error "MODEL_REGISTRY_TOKEN environment variable is not set"
    exit 1
fi

# Load environment variables from config file if it exists
if [[ -f "$PROJECT_DIR/tests/e2e_config.env" ]]; then
    print_status "Loading environment variables from e2e_config.env"
    source "$PROJECT_DIR/tests/e2e_config.env"
fi

# Display configuration
print_status "Configuration:"
echo "  Host: $MODEL_REGISTRY_HOST"
echo "  Port: ${MODEL_REGISTRY_PORT:-8080}"
echo "  Secure: ${MODEL_REGISTRY_SECURE:-false}"
echo "  Token: ${MODEL_REGISTRY_TOKEN:0:10}..." # Show first 10 chars for security

# Check if the package is installed
if ! python -c "import modelregistry_plugin" 2>/dev/null; then
    print_warning "modelregistry_plugin package not found. Installing..."
    cd "$PROJECT_DIR"
    uv build
    uv pip install dist/*.whl
fi

# Run the e2e tests
print_status "Running E2E tests..."
cd "$PROJECT_DIR"

# Run with verbose output and show local variables on failure
pytest tests/test_e2e.py -v -s --tb=short --showlocals

if [[ $? -eq 0 ]]; then
    print_success "All E2E tests passed!"
else
    print_error "Some E2E tests failed!"
    exit 1
fi 