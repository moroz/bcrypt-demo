#!/bin/sh -ex

export GO_ENV=prod

scripts/migrate_db.sh
./server
