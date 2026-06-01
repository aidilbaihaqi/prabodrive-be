#!/bin/bash

# ===========================================
# Migration Helper Script
# ===========================================

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Build database URL if not set
if [ -z "$DATABASE_URL" ]; then
    DATABASE_URL="postgres://${DB_USER:-postgres}:${DB_PASSWORD:-postgres}@${DB_HOST:-localhost}:${DB_PORT:-5432}/${DB_NAME:-myapp_db}?sslmode=${DB_SSL_MODE:-disable}"
fi

MIGRATIONS_PATH="./migrations"

usage() {
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  up          Run all pending migrations"
    echo "  down        Rollback the last migration"
    echo "  down-all    Rollback all migrations"
    echo "  version     Show current migration version"
    echo "  force N     Force set version to N"
    echo "  create NAME Create a new migration"
    echo ""
    echo "Examples:"
    echo "  $0 up"
    echo "  $0 down"
    echo "  $0 create create_products_table"
}

case "$1" in
    up)
        echo "Running migrations..."
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" up
        echo "✓ Migrations complete"
        ;;
    down)
        echo "Rolling back last migration..."
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" down 1
        echo "✓ Rollback complete"
        ;;
    down-all)
        echo "Rolling back all migrations..."
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" down -all
        echo "✓ All migrations rolled back"
        ;;
    version)
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" version
        ;;
    force)
        if [ -z "$2" ]; then
            echo "Error: Version number required"
            exit 1
        fi
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" force $2
        echo "✓ Version set to $2"
        ;;
    create)
        if [ -z "$2" ]; then
            echo "Error: Migration name required"
            exit 1
        fi
        migrate create -ext sql -dir $MIGRATIONS_PATH -seq $2
        echo "✓ Created migration: $2"
        ;;
    *)
        usage
        exit 1
        ;;
esac
