# Speedcubing Slovakia

Web app for regular online competitions in rubik's cube solving and a place for Slovak speedcubing community.

---

### Local development _(Linux)_

Make sure to have [Git](https://git-scm.com/download/linux) installed.

Clone this repository. `https://github.com/jakubdrobny/speedcubingslovakia.git`

---

#### Database

1. Install [PostgreSQL](https://www.postgresql.org/) with `sudo apt-get install postgres`.
2. Create some user, or configure password default user `postgres`.
3. Run `psql -U postgres` and then inside `psql` run `CREATE DATABASE <database_name>;` to create the database.
4. Exit `psql` with CTRL+D and inside `/backend/.env.developement` change `speedcubingslovakiadb_local` to your `<database_name>`, similarly for `username` and `password`.
5. Populate/initialize the database by running `psql -U <username> -d <database_name> -f initialize_db.sql` inside the `/database` directory.

#### Backend

1. Install [Go 1.22.0](https://go.dev/doc/install).
2. Run `go get <package_name>` for packages _models, constants, middlewares, main, controllers, utils, cube_ to install dependencies.
3. Add `SPEEDCUBINGSLOVAKIA_BACKEND_ENV=development` to your `~/.profile` and run `source ~/.profile` to realize the changes in current terminal.
4. Run `go run main/main.go` in the backend directory to start the server.

#### Scrambling

1. Install PM2 with `npm install -g pm2`.
2. Run `npm install` inside the `scrambling` directory to install dependencies.
3. Run `pm2 start index.js --name scrambling.service` inside the `scrambling` directory.

#### Frontend

1. Install `npm` with `apt install npm`.
2. Install latest [Node.js](https://nodejs.org/en) version with `npm i node@lts` and then all the dependencies with `npm install`.
3. Start the frontend with `npm start`.

### That's it :partying_face:

You should now be able to have a local instance of the app running :)
