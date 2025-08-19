#!/bin/sh

set -e

echo "===== STARTING DeletePastWCACompetitionsJob: $(date) ====="
/usr/bin/local/delete_past_wca_competitions_job >> /proc/1/fd/1 2>&1
echo "===== FINISHED DeletePastWCACompetitionsJob: $(date) ====="
