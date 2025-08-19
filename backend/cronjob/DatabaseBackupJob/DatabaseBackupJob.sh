#!/bin/sh

set -e

echo "===== STARTING DatabaseBackupJob: $(date) ====="
/usr/bin/local/database_backup_job>> /proc/1/fd/1 2>&1
echo "===== FINISHED database_backup_job: $(date) ====="
