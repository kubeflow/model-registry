#!/bin/bash

set -e

cd "$(dirname "$(readlink -f "$0")")/.."

if [ -z "$YQ" ]; then
  if [ -e "bin/yq" ]; then
    YQ="$(realpath "bin/yq")"
  else
    echo "Error: YQ is not set and bin/yq does not exist" >&2
    exit 1
  fi
fi

# Temporary files tracked for cleanup
TEMP_FILES=()

cleanup() {
    rm -f "${TEMP_FILES[@]}" 2>/dev/null || true
}
trap cleanup EXIT

# Register a temporary file for cleanup on exit
register_temp() {
    TEMP_FILES+=("$1")
}

usage() {
    echo "Usage: $0 [--check] <basename.yaml>"
    echo "  --check: Check for differences in the generated merged catalog specification."
    echo ""
    echo "This script merges the main catalog API with all plugin APIs to create"
    echo "a unified OpenAPI specification for documentation purposes."
    echo ""
    echo "Example: $0 catalog-spec.yaml"
    exit 0
}

# Load common schema names once (pipe-delimited string, e.g. "BaseResource|BaseResourceList|Error")
COMMON_SCHEMAS=""
if [[ -f "api/openapi/src/lib/common.yaml" ]]; then
    COMMON_SCHEMAS=$($YQ eval '.components.schemas | keys | join("|")' api/openapi/src/lib/common.yaml 2>/dev/null || echo "")
fi

# Check if a schema name is in the common schemas list.
is_common_schema() {
    local name="$1"
    [[ -n "$COMMON_SCHEMAS" && "|${COMMON_SCHEMAS}|" == *"|${name}|"* ]]
}

# Function to preprocess a plugin spec to avoid conflicts
preprocess_plugin_spec() {
    local spec_file="$1"
    local plugin_name="$2"
    local temp_file="$3"

    # Create capitalized prefix for schemas (e.g., mcp -> Mcp_)
    local schema_prefix="${plugin_name^}_"

    # Create lowercase prefix for operation IDs (e.g., mcp -> mcp_)
    local op_prefix="${plugin_name}_"

    # Start with original spec
    cp "$spec_file" "$temp_file"

    # 1. Prefix plugin-specific schema definitions (exclude common schemas)
    local schema_names
    schema_names=$($YQ eval '.components.schemas | keys | .[]' "$temp_file" 2>/dev/null || echo "")

    if [[ -n "$COMMON_SCHEMAS" ]]; then
        # Build a single yq expression to rename all plugin-specific schemas at once
        local yq_expr=""
        while IFS= read -r schema_name; do
            [[ -z "$schema_name" ]] && continue
            if ! is_common_schema "$schema_name"; then
                local new_name="${schema_prefix}${schema_name}"
                if [[ -n "$yq_expr" ]]; then
                    yq_expr="${yq_expr} | "
                fi
                yq_expr="${yq_expr}(.components.schemas[\"${new_name}\"] = .components.schemas[\"${schema_name}\"]) | del(.components.schemas[\"${schema_name}\"])"
            fi
        done <<< "$schema_names"

        if [[ -n "$yq_expr" ]]; then
            if ! $YQ eval -i "$yq_expr" "$temp_file" 2>/dev/null; then
                echo "Warning: failed to rename schemas for plugin $plugin_name" >&2
            fi
        fi
    else
        # Fallback: manually prefix non-BaseResource schemas
        local yq_expr=""
        while IFS= read -r schema_name; do
            if [[ -n "$schema_name" && "$schema_name" != "BaseResource" && "$schema_name" != "BaseResourceList" && "$schema_name" != "BaseResourceDates" ]]; then
                if [[ -n "$yq_expr" ]]; then
                    yq_expr="${yq_expr} | "
                fi
                yq_expr="${yq_expr}(.components.schemas[\"${schema_prefix}${schema_name}\"] = .components.schemas[\"${schema_name}\"]) | del(.components.schemas[\"${schema_name}\"])"
            fi
        done <<< "$schema_names"

        if [[ -n "$yq_expr" ]]; then
            if ! $YQ eval -i "$yq_expr" "$temp_file" 2>/dev/null; then
                echo "Warning: failed to rename schemas for plugin $plugin_name (fallback)" >&2
            fi
        fi
    fi

    # 2. Update all $ref pointers to use prefixed schema names (except for common schemas)
    # First, prefix all schema references throughout the document
    sed -i 's|#/components/schemas/\([A-Za-z][A-Za-z0-9]*\)|#/components/schemas/'"${schema_prefix}"'\1|g' "$temp_file"
    # Then, un-prefix each common schema individually
    if [[ -n "$COMMON_SCHEMAS" ]]; then
        IFS='|' read -ra SCHEMA_ARRAY <<< "$COMMON_SCHEMAS"
        for schema in "${SCHEMA_ARRAY[@]}"; do
            sed -i 's|#/components/schemas/'"${schema_prefix}${schema}"'|#/components/schemas/'"${schema}"'|g' "$temp_file"
        done
    else
        sed -i 's|#/components/schemas/'"${schema_prefix}"'BaseResource|#/components/schemas/BaseResource|g' "$temp_file"
        sed -i 's|#/components/schemas/'"${schema_prefix}"'BaseResourceList|#/components/schemas/BaseResourceList|g' "$temp_file"
        sed -i 's|#/components/schemas/'"${schema_prefix}"'BaseResourceDates|#/components/schemas/BaseResourceDates|g' "$temp_file"
    fi

    # 3. Resolve external references to common schemas before merging
    sed -i 's|lib/common.yaml#/components/schemas/|#/components/schemas/|g' "$temp_file"

    # 4. Prefix operation IDs to ensure uniqueness
    sed -i 's/operationId: \(.*\)/operationId: '"${op_prefix}"'\1/' "$temp_file"

    # 5. Convert relative paths to absolute paths using the plugin's server base URL
    local plugin_base_url
    plugin_base_url=$($YQ eval '.servers[0].url' "$temp_file" 2>/dev/null || echo "")
    if [[ -n "$plugin_base_url" && "$plugin_base_url" != "null" ]]; then
        local yq_path_expr=""
        local paths
        paths=$($YQ eval '.paths | keys | .[]' "$temp_file" 2>/dev/null || echo "")
        while IFS= read -r path; do
            if [[ -n "$path" && "$path" != "null" ]]; then
                local absolute_path="${plugin_base_url}${path}"
                if [[ -n "$yq_path_expr" ]]; then
                    yq_path_expr="${yq_path_expr} | "
                fi
                yq_path_expr="${yq_path_expr}(.paths[\"${absolute_path}\"] = .paths[\"${path}\"]) | del(.paths[\"${path}\"])"
            fi
        done <<< "$paths"

        if [[ -n "$yq_path_expr" ]]; then
            if ! $YQ eval -i "$yq_path_expr" "$temp_file" 2>/dev/null; then
                echo "Warning: failed to rewrite paths for plugin $plugin_name" >&2
            fi
        fi

        # Remove the server configuration since paths are now absolute
        $YQ eval -i 'del(.servers)' "$temp_file" 2>/dev/null || true
    fi

    # 6. Remove the info section to avoid overwriting main catalog's info during merge
    $YQ eval -i 'del(.info)' "$temp_file" 2>/dev/null || true
}

CHECK=false
BASENAME=""
while [[ $# -gt 0 ]]; do
    case "$1" in
        --check)
            CHECK=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            if [[ "${1#-}" != "$1" ]]; then
                echo "Unknown option: $1"
                usage
            fi
            if [[ "$BASENAME" != "" ]]; then
                usage
            fi

            BASENAME=$1
            shift
            ;;
    esac
done

if [[ "$BASENAME" == "" ]]; then
    usage
fi

BASENAME=$(basename "$BASENAME")
MAIN_CATALOG="api/openapi/catalog.yaml"

if [[ ! -f "$MAIN_CATALOG" ]]; then
    echo "Main catalog specification not found at $MAIN_CATALOG"
    exit 1
fi

OUT_FILE="api/openapi/$BASENAME"
if [[ "$CHECK" == "true" ]]; then
    OUT_FILE="$(mktemp -t modelregistry_catalog_spec_tempXXXXXX).yaml"
    register_temp "$OUT_FILE"
fi

# Auto-discover plugin OpenAPI specifications (exclude src subdirectory)
PLUGIN_SPECS=()
while IFS= read -r spec; do
    PLUGIN_SPECS+=("$spec")
done < <(find catalog/plugins/*/api/openapi -maxdepth 1 -name "openapi.yaml" -type f 2>/dev/null | sort || true)

echo "Merging catalog specifications..."
echo "  Main catalog: $MAIN_CATALOG"

# Start with the main catalog specification
cp "$MAIN_CATALOG" "$OUT_FILE"

# Track which plugins we're merging for enhanced description
PLUGIN_NAMES=()

# Process each plugin spec
for plugin_spec in "${PLUGIN_SPECS[@]}"; do
    # Extract plugin name from path (e.g., catalog/plugins/mcp/... -> mcp)
    plugin_name=${plugin_spec#catalog/plugins/}
    plugin_name=${plugin_name%%/*}

    if [[ -z "$plugin_name" ]]; then
        echo "Warning: Could not extract plugin name from $plugin_spec, skipping..."
        continue
    fi

    echo "  Plugin: $plugin_name ($plugin_spec)"
    PLUGIN_NAMES+=("$plugin_name")

    # Preprocess the plugin spec to avoid conflicts
    temp_preprocessed="$(mktemp -t "preprocessed_${plugin_name}_XXXXXX").yaml"
    register_temp "$temp_preprocessed"
    preprocess_plugin_spec "$plugin_spec" "$plugin_name" "$temp_preprocessed"

    # Merge the preprocessed plugin spec with the main spec
    temp_merged="$(mktemp -t merged_tempXXXXXX).yaml"
    register_temp "$temp_merged"
    $YQ eval-all '. as $item ireduce ({}; . * $item)' "$OUT_FILE" "$temp_preprocessed" > "$temp_merged"
    mv "$temp_merged" "$OUT_FILE"
done

# Update the main spec's info section to reflect the merged content
if [[ ${#PLUGIN_NAMES[@]} -gt 0 ]]; then
    plugin_list=$(IFS=', '; echo "${PLUGIN_NAMES[*]}")
    $YQ eval -i '.info.description = .info.description + "\n\nThis unified specification includes APIs from the following plugins: '"$plugin_list"'."' "$OUT_FILE"

    # Add a custom extension to track included plugins
    plugins_json="[$(printf '"%s",' "${PLUGIN_NAMES[@]}" | sed 's/,$//')]"
    $YQ eval -i '.["x-catalog-plugins"] = '"$plugins_json" "$OUT_FILE"
fi

# Re-order the keys in the generated file (following merge_openapi.sh pattern)
$YQ eval -i '
    {
        "openapi": .openapi,
        "info": .info,
        "servers": .servers,
        "paths": .paths,
        "components": .components,
        "security": .security,
        "tags": .tags,
        "x-catalog-plugins": .["x-catalog-plugins"]
    } |
        sort_keys(.paths) |
        sort_keys(.components.schemas) |
        sort_keys(.components.responses) |
        sort_keys(.components.parameters)
' "$OUT_FILE"

if [[ "$CHECK" == "true" ]]; then
    diff -u "api/openapi/$BASENAME" "$OUT_FILE"
    exit $?
fi

echo "Merged catalog specification generated: $OUT_FILE"
if [[ ${#PLUGIN_NAMES[@]} -gt 0 ]]; then
    echo "Included plugins: $(IFS=', '; echo "${PLUGIN_NAMES[*]}")"
else
    echo "No plugin specifications found - using main catalog only"
fi
