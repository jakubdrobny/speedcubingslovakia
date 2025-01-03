#!/bin/sh

set -o allexport
. ~/.profile
set +o allexport

cd $SPEEDCUBINGSLOVAKIA_PATH/backend

echo "===== START: `date` =====" >> $SPEEDCUBINGSLOVAKIA_CRONJOB_LOGFILE_PATH/UpcomingWCACompetitionsJob.txt 2>&1
go run cronjob/UpcomingWCACompetitionsJob/UpcomingWCACompetitionsJob.go >> $SPEEDCUBINGSLOVAKIA_CRONJOB_LOGFILE_PATH/UpcomingWCACompetitionsJob.txt 2>&1
echo "===== END =====" >> $SPEEDCUBINGSLOVAKIA_CRONJOB_LOGFILE_PATH/UpcomingWCACompetitionsJob.txt 2>&1
