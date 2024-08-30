#!/bin/sh

set -o allexport
. ~/.profile
set +o allexport
cd /home/jakubd/skola/Tvorba-internetovych-aplikacii/speedcubingslovakia/backend
go run cronjob/WeeklyCompetitionJob.go >> /home/jakubd/cronjob_tmp_output/speedcubingslovakia_output.txt 2>&1