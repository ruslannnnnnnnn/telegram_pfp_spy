#!/bin/bash
set -e

set -a
source .env
set +a

MIGRATIONS_DIR="./migrations"

CLICKHOUSE_URL="tcp://localhost:9000/${CLICKHOUSE_DB}?username=${CLICKHOUSE_USER}&password=${CLICKHOUSE_PASSWORD}"

goose -dir "$MIGRATIONS_DIR" clickhouse "$CLICKHOUSE_URL" up
