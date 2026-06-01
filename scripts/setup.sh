#!/bin/bash

# ===========================================
# Setup Script for Go Backend Template
# ===========================================

set -e

echo "🚀 Setting up Go Backend Template..."

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

echo "✓ Go $(go version | awk '{print $3}')"

# Copy .env if not exists
if [ ! -f .env ]; then
    if [ -f .env.example ]; then
        cp .env.example .env
        echo "✓ Created .env from .env.example"
    fi
fi

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
echo "✓ Dependencies downloaded"

# Install tools
echo "🔧 Installing development tools..."

# Install golang-migrate
if ! command -v migrate &> /dev/null; then
    echo "  Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# Install air for hot reload
if ! command -v air &> /dev/null; then
    echo "  Installing air (hot reload)..."
    go install github.com/air-verse/air@latest
fi

# Install golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    echo "  Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

echo "✓ Development tools installed"

# Verify setup
echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Edit .env with your configuration"
echo "  2. Start PostgreSQL and Redis"
echo "  3. Run 'make migrate-up' to run migrations"
echo "  4. Run 'make run' to start the server"
echo ""
echo "Or use Docker:"
echo "  make docker-up"
echo ""
