# create_migration:
# 	migrate create -ext=sql -dir=migrations -seq $(TITLE)
#
# migrate_up:
# 	migrate -path=migrations/ -database "${DB_URL}" -verbose up $(STEPS)
#
# migrate_down:
# 	migrate -path=migrations/ -database "${DB_URL}" -verbose down $(STEPS)
#
# migrate_force:
# 	migrate -path=migrations/ -database "${DB_URL}" force $(VERSION)

create_migration:
	@if [ -z "$(TITLE)" ]; then echo "Usage: make create_migration TITLE=your_migration_title"; exit 1; fi
	@echo "Creating migration: $(TITLE)"
	MIGRATE_CMD=create MIGRATE_TITLE=$(TITLE) docker compose -f $(COMPOSE_FILE_PATH) run --rm migrate

migrate_up:
	@echo "Migrating up..."
	MIGRATE_CMD=up MIGRATE_STEPS=$(STEPS) docker compose -f $(COMPOSE_FILE_PATH) run --rm migrate

migrate_down:
	@echo "Migrating down..."
	MIGRATE_CMD=down MIGRATE_STEPS=$(STEPS) docker compose -f $(COMPOSE_FILE_PATH) run --rm migrate

migrate_force:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make migrate_force VERSION=your_version_number"; exit 1; fi
	@echo "Forcing migration to version $(VERSION)..."
	MIGRATE_CMD=force MIGRATE_VERSION=$(VERSION) docker compose -f $(COMPOSE_FILE_PATH) run --rm migrate

.PHONY: create_migration migrate_up migrate_down migrate_force
