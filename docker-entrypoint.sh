#!/bin/sh
set -e

# Determine which mode to run
if [ "$APP_MODE" = "web" ]; then
    echo "Starting CodexGigantus Web GUI..."
    exec /app/codexgigantus-web "$@"
else
    echo "Starting CodexGigantus CLI..."
    exec /app/codexgigantus-cli "$@"
fi
