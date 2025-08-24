FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN \
  CGO_ENABLED=0 go build -o /app/bin/weekly_competition_job ./cronjob/WeeklyCompetitionJob/WeeklyCompetitionJob.go & \
  CGO_ENABLED=0 go build -o /app/bin/database_backup_job ./cronjob/DatabaseBackupJob/DatabaseBackupJob.go & \
  CGO_ENABLED=0 go build -o /app/bin/upcoming_wca_competitions_job ./cronjob/UpcomingWCACompetitionsJob/UpcomingWCACompetitionsJob.go & \
  CGO_ENABLED=0 go build -o /app/bin/delete_past_wca_competitions_job ./cronjob/DeletePastWCACompetitionsJob/DeletePastWCACompetitionsJob.go & \
  wait

FROM alpine:latest

RUN apk add --no-cache busybox-suid postgresql-client tzdata

WORKDIR /app

RUN mkdir -p jobs logs config db_backups

COPY --from=builder /app/bin/* /usr/local/bin

COPY ./cronjob/run-job.sh /app/jobs/run-job.sh

COPY ./cronjob/crontab /etc/crontabs/root

RUN chmod +x /app/jobs/run-job.sh
RUN chmod 0644 /etc/crontabs/root

CMD ["crond", "-f", "-l", "8"]
