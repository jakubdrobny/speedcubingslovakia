#!/bin/sh

set -e

if [ "$#" -ne 2 ]; then
  echo "Error: Invalid number of arguments." >&2
  echo "Usage: $0 <path_to_executable> <job_name_for_logging>" >&2
  exit 1
fi

EXECUTABLE_PATH=$1
JOB_NAME=$2

echo "===== STARTING ${JOB_NAME}: $(date) ====="
${EXECUTABLE_PATH} >> /proc/1/fd/1 2>&1
echo "===== FINISHED ${JOB_NAME}: $(date) ====="
