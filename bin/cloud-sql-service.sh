#!/usr/bin/env bash

# bash strict settings
set -o errexit # exit on errors
set -o nounset # exit on use of uninitialized variable
set -o pipefail

isRunning() {
  local STATE=`gcloud sql instances list --filter="NAME:${CLOUDSQL_INSTANCE_NAME}" --format='get(state)' --quiet --project="${GCP_PROJECT_ID}"`
  if [[ ${STATE} == 'RUNNABLE' ]]; then
    return 0 # bash true
  else
    return 1 # bash false
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
    gcloud sql instances patch "${CLOUDSQL_INSTANCE_NAME}" --activation-policy ALWAYS --project="${GCP_PROJECT_ID}" --quiet \
      || ( echo "Startup may be taking a little extra time. We'll give it another 5 minutes."; \
           gcloud beta sql operations wait --quiet $(gcloud sql operations list --instance="${CLOUDSQL_INSTANCE_NAME}" --filter='status=RUNNING' --format="value(NAME)" --project="${GCP_PROJECT_ID}") --project="${GCP_PROJECT_ID}" )
    ;;
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
