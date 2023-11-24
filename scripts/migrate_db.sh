#!/bin/sh -e

OPERATION="${@:-up}"

echo "Executing migrations with direction: $OPERATION"

if [ "$DATABASE_URL" = "" ]; then
    echo "FATAL: Environment variable DATABASE_URL is not set!"
    exit 1
fi

migrate -database $DATABASE_URL -path db/migrations $OPERATION
