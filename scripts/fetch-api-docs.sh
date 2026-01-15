#!/bin/bash
#
# Fetches the schack.se OpenAPI spec and manages versions
#
# Usage: ./scripts/fetch-api-docs.sh
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DOCS_DIR="$PROJECT_ROOT/api-specs"
API_URL="https://member.schack.se/memdb/v3/api-docs"
CURRENT_FILE="$DOCS_DIR/schack-se-api-docs.json"
OLD_FILE="$DOCS_DIR/schack-se-api-docs-old.json"
TEMP_FILE="$DOCS_DIR/.api-docs-temp.json"

# Create docs directory if it doesn't exist
mkdir -p "$DOCS_DIR"

echo "Fetching schack.se API docs from $API_URL..."

# Download to temp file
if ! curl -s "$API_URL" | jq '.' > "$TEMP_FILE"; then
    echo "Error: Failed to fetch or parse API docs"
    rm -f "$TEMP_FILE"
    exit 1
fi

# Check if we got valid JSON
if [ ! -s "$TEMP_FILE" ]; then
    echo "Error: Downloaded file is empty"
    rm -f "$TEMP_FILE"
    exit 1
fi

# Get version info from the spec
NEW_VERSION=$(jq -r '.info.version // "unknown"' "$TEMP_FILE")
NEW_TITLE=$(jq -r '.info.title // "unknown"' "$TEMP_FILE")
echo "Downloaded: $NEW_TITLE v$NEW_VERSION"

# Check if current file exists
if [ -f "$CURRENT_FILE" ]; then
    OLD_VERSION=$(jq -r '.info.version // "unknown"' "$CURRENT_FILE")
    echo "Existing version: v$OLD_VERSION"

    # Compare files (ignoring whitespace differences since we pretty-print)
    if diff -q <(jq -S '.' "$CURRENT_FILE") <(jq -S '.' "$TEMP_FILE") > /dev/null 2>&1; then
        echo "No changes detected in API spec"
        rm -f "$TEMP_FILE"
        exit 0
    fi

    echo ""
    echo "Changes detected! Comparing specs..."
    echo "=================================="

    # Show summary of changes
    OLD_PATHS=$(jq -r '.paths | keys | length' "$CURRENT_FILE")
    NEW_PATHS=$(jq -r '.paths | keys | length' "$TEMP_FILE")
    echo "Endpoints: $OLD_PATHS -> $NEW_PATHS"

    # List added paths
    ADDED_PATHS=$(comm -13 <(jq -r '.paths | keys | .[]' "$CURRENT_FILE" | sort) <(jq -r '.paths | keys | .[]' "$TEMP_FILE" | sort))
    if [ -n "$ADDED_PATHS" ]; then
        echo ""
        echo "Added endpoints:"
        echo "$ADDED_PATHS" | while read path; do echo "  + $path"; done
    fi

    # List removed paths
    REMOVED_PATHS=$(comm -23 <(jq -r '.paths | keys | .[]' "$CURRENT_FILE" | sort) <(jq -r '.paths | keys | .[]' "$TEMP_FILE" | sort))
    if [ -n "$REMOVED_PATHS" ]; then
        echo ""
        echo "Removed endpoints:"
        echo "$REMOVED_PATHS" | while read path; do echo "  - $path"; done
    fi

    # Show schemas changes
    OLD_SCHEMAS=$(jq -r '.components.schemas // .definitions | keys | length' "$CURRENT_FILE" 2>/dev/null || echo "0")
    NEW_SCHEMAS=$(jq -r '.components.schemas // .definitions | keys | length' "$TEMP_FILE" 2>/dev/null || echo "0")
    echo ""
    echo "Schemas: $OLD_SCHEMAS -> $NEW_SCHEMAS"

    echo ""
    echo "=================================="

    # Backup old file
    echo "Backing up old version to schack-se-api-docs-old.json"
    mv "$CURRENT_FILE" "$OLD_FILE"
fi

# Move temp file to current
mv "$TEMP_FILE" "$CURRENT_FILE"

echo ""
echo "API docs saved to: $CURRENT_FILE"
echo ""
echo "Quick stats:"
echo "  - OpenAPI version: $(jq -r '.openapi // .swagger' "$CURRENT_FILE")"
echo "  - API version: $NEW_VERSION"
echo "  - Endpoints: $(jq -r '.paths | keys | length' "$CURRENT_FILE")"
echo "  - Schemas: $(jq -r '.components.schemas // .definitions | keys | length' "$CURRENT_FILE" 2>/dev/null || echo "N/A")"