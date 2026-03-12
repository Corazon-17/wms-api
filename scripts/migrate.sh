#!/bin/bash

set -e

DB_URL=$1
MIGRATION_PATH=$2

echo "Acquiring migration lock..."

psql $DB_URL -c "SELECT pg_advisory_lock(123456);"

echo "Running migrations..."

migrate -path $MIGRATION_PATH -database "$DB_URL" drop -f

migrate -path $MIGRATION_PATH -database "$DB_URL" up

echo "Releasing migration lock..."

psql $DB_URL -c "SELECT pg_advisory_unlock(123456);"

echo "Migration complete."