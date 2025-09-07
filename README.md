# Speedcubing Slovakia

Web app for regular online competitions in Rubik's cube solving and a place for Slovak speedcubing community.

---

### Prerequisites

Make sure to have [Git](https://git-scm.com/download/linux), [docker](https://docs.docker.com/engine/install/ubuntu/) (and `docker compose`) installed.

### Running the application

1. Copy the `.env.example` file into a new `.env.development` file in the project root and fill in the environment variables:
    - `WCA_CLIENT_ID` and `WCA_CLIENT_SECRET` - go to `your WCA profile > Manage your applications > Create` and set the `name` to anything you like, `redirect uri` to `http://localhost:3000/login` and `scope` to `public+email` and then copy the created `client id` and `client secret` to the variables
    - `JWT_SECRET_KEY` - could be anything for local development
    - `MAIL_USERNAME` - email address from which to send the newsletter emails from and to which to send alerts about suspicous results
    - `MAIL_PASSWORD` - for gmail it has to be the [app password](https://support.google.com/accounts/answer/185833?hl=en)
    - `DRIVE_BACKUP_FOLDER_ID` - you can find this id in the url when you open the google drive directory in your browser
    - the paths in the variables should not be changed, since they are paths inside the docker containers, not your machine

2. Start the entire application (frontend, backend, database, scrambling service, cron jobs and the monitoring stack) with:

`docker compose -f docker-compose.dev.yml up -d --build`

You should have the entire app up and running :D
