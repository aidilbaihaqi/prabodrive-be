#!/bin/bash
set -e

# Load .env if present
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

DB_URL="postgres://${DB_USER:-postgres}:${DB_PASSWORD:-secret}@${DB_HOST:-localhost}:${DB_PORT:-5432}/${DB_NAME:-prabodrive}?sslmode=${DB_SSLMODE:-disable}"
MIGRATIONS_PATH="./migrations"

usage() {
    echo "Usage: $0 <command>"
    echo ""
    echo "Commands:"
    echo "  up           Run all pending migrations"
    echo "  down         Rollback 1 migration"
    echo "  down-all     Rollback semua migration"
    echo "  fresh        Rollback semua lalu migrate up (hapus semua data)"
    echo "  version      Cek versi migration saat ini"
    echo "  force N      Paksa set versi ke N (setelah error)"
    echo "  create NAME  Buat file migration baru"
}

case "$1" in
    up)
        echo "▶ Running migrations..."
        migrate -path $MIGRATIONS_PATH -database "$DB_URL" up
        echo "✓ Done"
        ;;
    down)
        echo "▶ Rolling back 1 migration..."
        migrate -path $MIGRATIONS_PATH -database "$DB_URL" down 1
        echo "✓ Done"
        ;;
    down-all)
        echo "▶ Rolling back ALL migrations..."
        migrate -path $MIGRATIONS_PATH -database "$DB_URL" down -all
        echo "✓ Done"
        ;;
    fresh)
        echo "▶ Dropping all tables and re-migrating..."
        migrate -path $MIGRATIONS_PATH -database "$DB_URL" down -all
        migrate -path $MIGRATIONS_PATH -database "$DB_URL" up
        echo "✓ Fresh migration complete"
        ;;
    version)
        migrate -path $MIGRATIONS_PATH -database "$DB_URL" version
        ;;
    force)
        [ -z "$2" ] && echo "Error: version number required" && exit 1
        migrate -path $MIGRATIONS_PATH -database "$DB_URL" force $2
        echo "✓ Version forced to $2"
        ;;
    create)
        [ -z "$2" ] && echo "Error: migration name required" && exit 1
        migrate create -ext sql -dir $MIGRATIONS_PATH -seq $2
        echo "✓ Created: migrations/*_$2.{up,down}.sql"
        ;;
    *)
        usage
        exit 1
        ;;
esac
