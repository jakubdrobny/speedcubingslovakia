# Speedcubing Slovakia

Web app for regular online competitions in rubik's cube solving and a place for Slovak speedcubing community.

---

### Local development _(Linux)_

Make sure to have [Git](https://git-scm.com/download/linux) installed.

Clone this repository. `https://github.com/jakubdrobny/speedcubingslovakia.git`

---

#### Prerequisites

1. Install [PostgreSQL](https://www.postgresql.org/) with `sudo apt-get install postgres`.
2. Install [Go 1.22.0](https://go.dev/doc/install).
3. Install latest [Node.js with npm](https://nodejs.org/en/download).
4. Install PM2 with `npm install -g pm2`.

#### Database

1. Create some user, or configure password default user `postgres`.
2. Run `psql -U postgres` and then inside `psql` run `CREATE DATABASE <database_name>;` to create the database.
3. Exit `psql` with CTRL+D and create the `backend/.env.developement` file similarly to [`backend/.env.development.example`](./backend/.env.development.example), and change the `<username>`, `<password>`, and `<database_name>` placeholders to your owns (e.g.: `DB_URL=postgres://postgres:verySecurePw12345@localhost:5432/speedcubingslovakia_example`).
4. ~~Populate/initialize the database by running `psql -U <username> -d <database_name> -f initialize_db.sql` inside the `database` directory.~~ Populate/initialize the database by running `make migrate_up` inside the `database` directory.

#### Backend

1. Run `go mod tidy` to install dependencies.
2. Add `export SPEEDCUBINGSLOVAKIA_BACKEND_ENV=development` to your `~/.profile` and run `source ~/.profile` to realize the changes in current terminal.
3. Run `go run main/main.go` in the backend directory to start the server.

#### Scrambling

1. Run `npm install` inside the `scrambling` directory to install dependencies.
2. Run `npm run start_service` inside the `scrambling` directory.

#### Frontend

1. Run `npm install` inside the `frontend` directory to install dependencies.
2. Start the frontend with `npm run dev`.

### That's it :partying_face:

You should now be able to have a local instance of the app running. :)
