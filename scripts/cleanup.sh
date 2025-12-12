#!/usr/bin/env bash

set -e

MR_NAMESPACE="${MR_NAMESPACE:-kubeflow}"
TEST_DB_NAME="${TEST_DB_NAME:-metadb}"
MYSQL_USER_NAME="${MYSQL_USER_NAME:-root}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-test}"
POSTGRES_USER="${POSTGRES_USER:-root}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-test}"
DEPLOY_MANIFEST_DB="${DEPLOY_MANIFEST_DB:-db}" # subdirectory of manifests/kustomize/overlays to select which database: 'db' (MySQL) or 'postgres'

# transaction start commands are different between sqlite and mysql
PARTIAL_SQL_CMD=$(
    cat <<EOF
DELETE FROM Artifact;
DELETE FROM ArtifactProperty;
DELETE FROM Association;
DELETE FROM Attribution;
DELETE FROM Context;
DELETE FROM ContextProperty;
DELETE FROM Event;
DELETE FROM EventPath;
DELETE FROM Execution;
DELETE FROM ExecutionProperty;
DELETE FROM ParentContext;
COMMIT;
EOF
)
POSTGRES_PARTIAL_SQL_CMD=$(
    cat <<EOF
DELETE FROM "Artifact";
DELETE FROM "ArtifactProperty";
DELETE FROM "Association";
DELETE FROM "Attribution";
DELETE FROM "Context";
DELETE FROM "ContextProperty";
DELETE FROM "Event";
DELETE FROM "EventPath";
DELETE FROM "Execution";
DELETE FROM "ExecutionProperty";
DELETE FROM "ParentContext";
COMMIT;
EOF
)

if [[ -n "$LOCAL" ]]; then
    echo 'Cleaning up local sqlite DB'

    sqlite3 test/config/ml-metadata/metadata.sqlite.db <<<"BEGIN TRANSACTION; $PARTIAL_SQL_CMD"
elif [[ "$DEPLOY_MANIFEST_DB" == "postgres" ]]; then
    echo -n 'Cleaning up kubernetes PostgreSQL DB...'

    kubectl exec -n "$MR_NAMESPACE" \
        "$(kubectl get pods -l component=db -o jsonpath="{.items[0].metadata.name}" -n "$MR_NAMESPACE")" \
        -- psql -U "$POSTGRES_USER" -d "$TEST_DB_NAME" -c "BEGIN; $POSTGRES_PARTIAL_SQL_CMD"

    echo -n 'Done cleaning up kubernetes PostgreSQL DB'
else
    echo -n 'Cleaning up kubernetes MySQL DB...'

    kubectl exec -n "$MR_NAMESPACE" \
        "$(kubectl get pods -l component=db -o jsonpath="{.items[0].metadata.name}" -n "$MR_NAMESPACE")" \
        -- mysql -h 127.0.0.1 -u "$MYSQL_USER_NAME" -p"$MYSQL_ROOT_PASSWORD" -D "$TEST_DB_NAME" -e "START TRANSACTION; $PARTIAL_SQL_CMD; COMMIT;"

    echo -n 'Done cleaning up kubernetes MySQL DB'
fi
