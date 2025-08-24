#!/bin/sh

set -e

if [ "$#" -ne 2 ]; then
  echo "Error: Invalid number of arguments." 2>&1
  echo "Usage: $0 <path_to_executable> <job_name_for_logging>" 2>&1
  exit 1
fi

EXECUTABLE_PATH=$1
JOB_NAME=$2
LOGS_PATH="/app/logs/${JOB_NAME}.txt"

echo "===== STARTING ${JOB_NAME}: $(date) =====" >> "${LOGS_PATH}" 2>&1
${EXECUTABLE_PATH} >> "${LOGS_PATH}" 2>&1
echo "===== FINISHED ${JOB_NAME}: $(date) =====" >> "${LOGS_PATH}" 2>&1
