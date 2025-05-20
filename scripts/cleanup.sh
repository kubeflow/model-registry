#!/usr/bin/env bash

set -e

MR_NAMESPACE="${MR_NAMESPACE:-kubeflow}"
TEST_DB_NAME="${TEST_DB_NAME:-metadb}"
MYSQL_USER_NAME="${MYSQL_USER_NAME:-root}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-test}"

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

if [[ -n "$LOCAL" ]]; then
    echo 'Cleaning up local sqlite DB'

    sqlite3 test/config/ml-metadata/metadata.sqlite.db <<<"BEGIN TRANSACTION; $PARTIAL_SQL_CMD"
else
    echo -n 'Cleaning up kubernetes MySQL DB...'

    kubectl exec -n "$MR_NAMESPACE" \
        "$(kubectl get pods -l component=db -o jsonpath="{.items[0].metadata.name}" -n "$MR_NAMESPACE")" \
        -- mysql -u "$MYSQL_USER_NAME" -p"$MYSQL_ROOT_PASSWORD" -D "$TEST_DB_NAME" -e "START TRANSACTION; $PARTIAL_SQL_CMD; COMMIT;"

    echo -n 'Done cleaning up kubernetes MySQL DB'
fi
