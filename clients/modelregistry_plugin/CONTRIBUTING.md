# Contributing to Model Registry MLflow Plugin

Thank you for your interest in contributing to the Model Registry MLflow Plugin! This guide is for developers who want to contribute to the plugin's development, testing, and maintenance.

## Table of Contents

- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Style and Standards](#code-style-and-standards)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Development Setup

### Prerequisites

- Python 3.9+
- Go 1.21+ (for building Model Registry server)
- Docker (for local e2e tests)
- `uv` (recommended) or `pip` for Python dependency management

### Initial Setup

```bash
# Clone the repository
git clone <repository-url>
cd model-registry/clients/modelregistry_plugin

# Install dependencies
uv sync

# Install in development mode
uv pip install -e .

# Verify setup
make verify-entry-point
```

### Development Environment

```bash
# Activate virtual environment (if using uv)
source .venv/bin/activate

# Install pre-commit hooks
pre-commit install
```

## Project Structure

```
clients/modelregistry_plugin/
├── modelregistry_plugin/          # Main plugin package
│   ├── __init__.py
│   ├── auth.py                   # Authentication handling
│   ├── store.py                  # Main ModelRegistryStore implementation
│   └── utils.py                  # Utility functions
├── tests/                        # Test suite
│   ├── __init__.py              # Test package initialization
│   ├── test_store.py            # Unit tests
│   ├── test_e2e.py              # Remote e2e tests
│   ├── test_e2e_local.py        # Local e2e tests
│   ├── test_auth.py             # Authentication tests
│   └── test_utils.py            # Utility tests
├── docs/                         # Documentation
├── pyproject.toml               # Project configuration
├── Makefile                     # Development tasks
└── README.md                    # End-user documentation
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Follow the [Code Style](#code-style-and-standards) guidelines
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run specific test categories
make test-e2e-local    # Local e2e tests (recommended)
make test-e2e          # Remote e2e tests (optional)

# Run linting and formatting
make lint
make format
```

### 4. Build and Verify

```bash
# Build the package
make build

# Install and verify entry point
make install-dev
make verify-entry-point
```

## Testing

### Test Categories

#### Unit Tests (`make test`)
- Fast, isolated tests for individual functions and classes
- No external dependencies
- Run on every commit and PR

#### Local E2E Tests (`make test-e2e-local`)
- **Recommended for most development**
- Starts local MLMD server (Docker) and Model Registry server
- Tests full plugin integration with real server
- No external dependencies or credentials required

#### Remote E2E Tests (`make test-e2e`)
- **Advanced integration testing against real environments**
- Requires remote Model Registry server and authentication
- Tests against production-like environments
- Useful for release validation and advanced debugging

### Running Tests

```bash
# Quick development cycle
make test                    # Unit tests with coverage (e2e tests excluded)
make test-e2e-local         # Local e2e tests
make test-e2e              # Remote e2e tests (requires MLFLOW_TRACKING_URI)

# Full test suite
make test && make test-e2e-local

# Run all tests including e2e (requires MLFLOW_TRACKING_URI to be set)
export MLFLOW_TRACKING_URI="modelregistry://your-host:port"
uv run pytest tests/ -v -s

# Specific test file
uv run pytest tests/test_store.py -k "test_search_runs"

# Debug mode
uv run pytest tests/ -v -s --log-cli-level=DEBUG -k "not e2e"
```

### Remote E2E Test Setup

Remote e2e tests require access to a running Model Registry server. Follow these steps to set them up:

#### Prerequisites

1. **Model Registry Server Access**
   - Local development server
   - Kubernetes cluster deployment
   - Cloud-hosted instance

2. **Authentication Token**
   - Valid token with sufficient permissions
   - Permissions for experiments, runs, metrics, and artifacts

3. **Network Connectivity**
   ```bash
   # Test basic connectivity
   curl -H "Authorization: Bearer YOUR_TOKEN" \
        "http://YOUR_HOST:PORT/health"
   ```

#### Configuration

**Required Environment Variables**
```bash
# Required: MLflow tracking URI for the Model Registry server
export MLFLOW_TRACKING_URI="modelregistry://your-host:port"

# Optional: Authentication token (if required)
export MODEL_REGISTRY_TOKEN="your-auth-token"
```

**Example Configuration**
```bash
# For HTTP server
export MLFLOW_TRACKING_URI="modelregistry://model-registry.example.com:8080"

# For HTTPS server
export MLFLOW_TRACKING_URI="modelregistry+https://model-registry.example.com:443"
export MODEL_REGISTRY_TOKEN="your-auth-token"
```

#### Running Remote E2E Tests

**Using Makefile (Recommended)**
```bash
# Set the tracking URI
export MLFLOW_TRACKING_URI="modelregistry+https://your-host:port"
export MODEL_REGISTRY_TOKEN="your-auth-token"

# Run e2e tests
make test-e2e
```

**Manual Execution**
```bash
# Set environment variables
export MLFLOW_TRACKING_URI="modelregistry+https://your-host:port"
export MODEL_REGISTRY_TOKEN="your-auth-token"

# Run e2e tests
uv run pytest tests/test_e2e.py -v -s
```

#### Test Coverage

Remote e2e tests cover:
- ✅ Connection and authentication
- ✅ Experiment management (create, retrieve, search, lifecycle)
- ✅ Run management (create, update, lifecycle)
- ✅ Data logging (metrics, parameters, tags)
- ✅ Search and filtering with pagination
- ✅ MLflow integration and entry point registration

#### Troubleshooting Remote E2E Tests

**Common Issues:**

1. **MLFLOW_TRACKING_URI Not Set**
   ```bash
   # Error: MLFLOW_TRACKING_URI is not set
   # Solution: Set the tracking URI
   export MLFLOW_TRACKING_URI="modelregistry://your-host:port"
   ```

2. **Connection Refused**
   ```bash
   # Verify server is running and accessible
   curl http://HOST:PORT/health
   ```

3. **Authentication Failed**
   ```bash
   # Test token validity
   curl -H "Authorization: Bearer TOKEN" http://HOST:PORT/health
   ```

4. **SSL/TLS Errors**
   ```bash
   # For HTTPS servers, ensure proper certificate configuration
   # The modelregistry:// URI scheme handles HTTPS automatically
   ```

5. **Entry Point Not Found**
   ```bash
   # Rebuild and reinstall
   uv build && uv pip install dist/*.whl --force-reinstall
   ```

**Debug Mode:**
```bash
# Enable debug logging
export LOG_LEVEL=DEBUG

# Run with debug output
uv run pytest tests/test_e2e.py -v -s --log-cli-level=DEBUG
```

### Test Best Practices

1. **Write tests for new functionality**
2. **Use descriptive test names**
3. **Test both success and failure cases**
4. **Mock external dependencies in unit tests**
5. **Use fixtures for common setup**
6. **Clean up test data after tests**
7. **Use unique names (UUIDs/timestamps) for e2e tests**
8. **Monitor test duration and optimize slow tests**

## Code Style and Standards

### Python Code Style

- **Formatting**: Use `ruff` for formatting and linting
- **Type hints**: Use type hints for all function parameters and return values
- **Docstrings**: Use Google-style docstrings for public functions and classes
- **Imports**: Use absolute imports, organize imports with `isort`

### Code Quality

```bash
# Format code
make format

# Check code quality
make lint

# Type checking (if mypy is configured)
uv run mypy modelregistry_plugin/
```

### Commit Message Format

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Examples:
- `feat(store): add support for batch metric logging`
- `fix(auth): handle expired token gracefully`
- `test(e2e): add local server test suite`
- `docs(readme): update installation instructions`

## Submitting Changes

### 1. Prepare Your Changes

```bash
# Ensure all tests pass
make test && make test-e2e-local

# Check code quality
make lint
make format

# Build and verify
make build
make verify-entry-point
```

### 2. Create a Pull Request

1. Push your branch to the repository
2. Create a pull request with:
   - Clear description of changes
   - Link to related issues
   - Test results summary
   - Screenshots (if UI changes)

### 3. PR Review Process

- All PRs require at least one review
- CI must pass (unit tests + local e2e tests)
- Code must follow style guidelines
- Documentation must be updated

## Release Process

### Version Management

The project uses semantic versioning (MAJOR.MINOR.PATCH):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Steps

1. **Update version** in `pyproject.toml`
2. **Update changelog** (if maintained)
3. **Run full test suite**:
   ```bash
   make test
   make test-e2e-local
   make test-e2e  # If remote server available
   ```
4. **Build and test package**:
   ```bash
   make build
   make install-dev
   make verify-entry-point
   ```
5. **Create release tag**:
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```
6. **Publish to PyPI** (if applicable)

### Pre-release Testing

Before each release:

1. **Run all test suites**
2. **Test with different Python versions** (if supported)
3. **Verify MLflow compatibility** with target MLflow version
4. **Test installation** in clean environment
5. **Verify entry point registration**

## Troubleshooting

### Common Issues

#### Build Failures
```bash
# Clean and rebuild
make clean
make build
```

#### Test Failures
```bash
# Check if Model Registry server is running (for e2e tests)
# Verify your MLFLOW_TRACKING_URI is set correctly
echo $MLFLOW_TRACKING_URI

# Run with debug output
uv run pytest -v -s --log-cli-level=DEBUG
```

#### Entry Point Issues
```bash
# Reinstall package
make install-dev
make verify-entry-point
```

#### Docker Issues (for local e2e tests)
```bash
# Check Docker is running
docker ps

# Clean up containers
docker system prune -f
```

#### Remote E2E Test Issues
```bash
# Verify MLFLOW_TRACKING_URI is set
echo $MLFLOW_TRACKING_URI

# Check connectivity (extract host:port from URI)
# For modelregistry://host:port, test with:
curl http://HOST:PORT/health

# Test with debug logging
LOG_LEVEL=DEBUG uv run pytest tests/test_e2e.py -v -s
```

### Getting Help

- **Issues**: Create an issue with detailed description and reproduction steps
- **Discussions**: Use GitHub Discussions for questions and ideas
- **Documentation**: Check existing docs in `docs/` directory

## Development Tools

### Useful Makefile Targets

```bash
make help              # Show all available targets
make dev               # Complete development setup
make test-e2e-local    # Run local e2e tests
make test-e2e          # Run remote e2e tests
make test-e2e-setup    # Test e2e setup
make show-stores       # Show available MLflow tracking stores
make verify-entry-point # Verify plugin registration
```

### IDE Configuration

#### VS Code
- Install Python extension
- Configure `ruff` for formatting and linting
- Set up debugging configuration for tests

#### PyCharm
- Configure project interpreter to use virtual environment
- Set up code style to match project standards
- Configure test runner for pytest

## Contributing Guidelines

### What to Contribute

- **Bug fixes**: Reproduce the issue, fix it, add tests
- **New features**: Discuss in issues first, implement with tests
- **Documentation**: Improve clarity, add examples, fix typos
- **Tests**: Add missing test coverage, improve existing tests
- **Performance**: Profile and optimize slow operations

### What Not to Contribute

- **Breaking changes** without discussion and approval
- **Features** that don't align with project goals
- **Code** without tests
- **Changes** that break existing functionality

### Code of Conduct

- Be respectful and inclusive
- Focus on technical merit
- Help others learn and contribute
- Follow project conventions and standards

---

Thank you for contributing to the Model Registry MLflow Plugin! Your contributions help make this project better for everyone. 