create_migration:
	migrate create -ext=sql -dir=migrations -seq $(TITLE)

migrate_up:
	migrate -path=migrations/ -database "${DB_URL}" -verbose up $(STEPS)

migrate_down:
	migrate -path=migrations/ -database "${DB_URL}" -verbose down $(STEPS)

migrate_force:
	migrate -path=migrations/ -database "${DB_URL}" force $(VERSION)
