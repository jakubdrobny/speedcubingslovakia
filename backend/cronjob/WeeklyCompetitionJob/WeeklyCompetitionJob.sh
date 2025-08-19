#!/bin/sh

set -e

echo "===== STARTING WeeklyCompetitionJob: $(date) ====="
/usr/bin/local/weekly_competition_job >> /proc/1/fd/1 2>&1
echo "===== FINISHED WeeklyCompetitionJob: $(date) ====="
