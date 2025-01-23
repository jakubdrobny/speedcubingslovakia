include ./.env.$(SPEEDCUBINGSLOVAKIA_BACKEND_ENV)
export

create_migration:
	migrate create -ext=sql -dir=migrations -seq $(TITLE)

migrate_up:
	migrate -path=migrations/ -database "${DB_URL}" -verbose up

migrate_up_steps:
	migrate -path=migrations/ -database "${DB_URL}" -verbose up $(STEPS)

migrate_down:
	migrate -path=migrations/ -database "${DB_URL}" -verbose down

migrate_down_steps:
	migrate -path=migrations/ -database "${DB_URL}" -verbose down $(STEPS)

migrate_force:
	migrate -path=migrations/ -database "${DB_URL}" force $(VERSION)
