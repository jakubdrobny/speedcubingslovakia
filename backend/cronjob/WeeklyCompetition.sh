#!/bin/sh

set -o allexport
. ~/.profile
set +o allexport
cd ~/skola/Tvorba-internetovych-aplikacii/speedcubingslovakia/backend
go run cronjob/WeeklyCompetitionJob.go >> ~/cronjob_tmp_output/WeeklyCompetitionJob.txt 2>&1