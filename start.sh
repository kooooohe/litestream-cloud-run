#!/bin/sh
set -e

if [ ! -f /app/db ]; then
	litestream restore -if-replica-exists -o /app/db "${BACKUP_URL}"
fi

exec litestream replicate -exec "/app/main -dp /app/db"
