#!/usr/bin/env bash

# bash strict settings
set -o errexit # exit on errors
set -o nounset # exit on use of uninitialized variable
set -o pipefail

isRunning() {
  local STATE=`gcloud sql instances list --filter="NAME:${CLOUDSQL_INSTANCE_NAME}" --format='get(state)' --quiet --project="${GCP_PROJECT_ID}"`
  if [[ ${STATE} == 'RUNNABLE' ]]; then
    return 0
  else
    return 1
  fi
}

ACTION="${1-}"

case "$ACTION" in
  name)
    echo "cloud-sql-service";;
  myorder)
    echo 0;;
  status)
    if isRunning; then
      echo "running"
    else
      echo "stopped"
    fi;;
  start)
    gcloud sql instances patch "${CLOUDSQL_INSTANCE_NAME}" --activation-policy ALWAYS --project="${GCP_PROJECT_ID}" --quiet;;
  stop)
    gcloud sql instances patch "${CLOUDSQL_INSTANCE_NAME}" --activation-policy NEVER --project="${GCP_PROJECT_ID}" --quiet;;
  restart)
    gcloud sql instances restart "${CLOUDSQL_INSTANCE_NAME}" --project="${GCP_PROJECT_ID}" --quiet;;
  param-default)
    echo '';;
  *)
    # TODO: library-ize and use 'echoerrandexit'
    echo "Unknown action '${ACTION}'." >&2
    exit 1;;
esac
