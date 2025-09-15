#!/bin/sh
set -e

CMD="${MIGRATE_CMD:-up}"
STEPS="${MIGRATE_STEPS}"
VERSION="${MIGRATE_VERSION}"
TITLE="${MIGRATE_TITLE}"
DATABASE_URL="${DB_URL}"

MIGRATE_BASE_CMD="migrate -path=/migrations -database ${DATABASE_URL} -verbose"

if [ "$CMD" = "create" ]; then
  if [ -z "$TITLE"]; then
    echo "Error: MIGRATE_TITLE is require for the 'create' command."
    exit 1
  fi

  echo "Creating migration with title: $TITLE"
  migrate create -ext sql -dir /migrations -seq "$TITLE"
  exit 0
fi

if [ -z "$DATABASE_URL" ]; then
  echo "Error: DB_URL environment variable not set."
  exit 1
fi

case "$CMD" in
  up)
    echo "Running migrations: up"
    if [ -n "$STEPS" ]; then
      $MIGRATE_BASE_CMD up "$STEPS"
    else
      $MIGRATE_BASE_CMD up
    fi
    ;;
  down)
    echo "Running migrations: down"
    if [ -n "$STEPS" ]; then
      $MIGRATE_BASE_CMD down "$STEPS"
    else
      $MIGRATE_BASE_CMD down
    fi
    ;;
  force)
    if [ -z "$VERSION" ]; then
      echo "Error: MIGRATE_VERSION is required for the 'force' command."
      exit 1
    fi
    echo "Forcing migration to version: $VERSION"
    $MIGRATE_BASE_CMD force "$VERSION"
    ;;
  *)
    echo "Error: Unknown command '$CMD'"
    exit 1
    ;;
esac 
