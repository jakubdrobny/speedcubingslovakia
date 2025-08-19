#!/bin/sh

set -e

echo "===== STARTING UpcomingWCACompetitionsJob: $(date) ====="
/usr/bin/local/upcoming_wca_competitions_job >> /proc/1/fd/1 2>&1
echo "===== FINISHED UpcomingWCACompetitionsJob: $(date) ====="
