# Speedcubing Slovakia

Web app for regular online competitions in rubik's cube solving and a place for Slovak speedcubing community.

## Local development *(Linux)*

Make sure to have [https://git-scm.com/download/linux]{Git} installed.

Clone this repository. `https://github.com/jakubdrobny/speedcubingslovakia.git`

#### Frontend

1. Install `npm` with `apt install npm`.
2. Install latest [https://nodejs.org/en]{Node.js} version with `npm i node@lts` and then all the dependencies with `npm install`.
3. Start the frontend with `npm start`.

#### Backend

1. Install [https://go.dev/doc/install]{Go 1.22.0}.
2. Run `go get <package_name>` for packages *models, constants, middlewares, main, controllers, utils* to install dependencies.
3. Run `go run main/main.go` in the backend directory to start the server.

#### Database
1. Install [https://www.postgresql.org/]{PostgreSQL} with `sudo apt-get install postgres`.
2. Configure postgres user password and 

