#!/bin/sh
# Apply pending DB migrations, then start the server. dbmate reads DATABASE_URL
# from the environment (set in docker-compose). --no-dump-schema avoids writing
# to the read-only image filesystem.
set -e

echo "Applying database migrations..."
dbmate --no-dump-schema --migrations-dir ./migrations up

echo "Starting server..."
exec ./jafa-server
