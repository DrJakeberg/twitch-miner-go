#!/usr/bin/env bash
set -euo pipefail

# Open the Config Editor
# Usage: ./tools/edit-config.sh [--config DIR] [--port PORT] [--tui] [--no-browser]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BINARY="$PROJECT_DIR/config-editor"

echo "Building config-editor..."
cd "$PROJECT_DIR"
go build -o config-editor ./cmd/config-editor

exec "$BINARY" "$@"
