#!/usr/bin/env bash

# bash strict settings
set -o errexit # exit on errors
set -o nounset # exit on use of uninitialized variable
set -o pipefail

ACTION="${1-}"; shift || true

# PORT=3306 # for mysql
PORT=5432 # for postgres

# 'cloud_sql_proxy' is essentially a wrapper that spawns a second process; but
# killing the parent does not kill the child, so we have to fallback to process
# grepping
isRunning() {
  set +e # grep exits with error if no match
  local PROC_COUNT=$(ps aux | grep cloud_sql_proxy | grep -v grep | wc -l)
  set -e
  # if [[ -f "$PID_FILE" ]] && kill -0 $(cat $PID_FILE) > /dev/null 2> /dev/null; then
  if (( $PROC_COUNT == 0 )); then
    return 1
  else
    return 0
  fi
}

startProxy() {
  local CLOUDSQL_CREDS="$HOME/.catalyst/creds/${CLOUDSQL_SERVICE_ACCT}.json"
  # We were using the following to capture the pid, but see note on 'isRunning'
  # bash -c "cd '${BASE_DIR}'; ( npx --no-install cloud_sql_proxy -instances='${CLOUDSQL_CONNECTION_NAME}'=tcp:$PORT -credential_file='${CLOUDSQL_CREDS}' & echo \$! >&3 ) 3> '${PID_FILE}' 2> '${SERV_LOG}' &"
  # Annoyingly, cloud_sql_proxy (at time of note) emits all logs to stderr.
  bash -c "cd '${BASE_DIR}'; npx --no-install cloud_sql_proxy -instances='${CLOUDSQL_PROXY_CONNECTION_NAME}'=tcp:$PORT -credential_file='$CLOUDSQL_CREDS' 2> '${SERV_LOG}' &"
}

stopProxy() {
  # See note in 'start'
  # kill $(cat "${PID_FILE}") && rm "${PID_FILE}"
  kill $(ps aux | grep cloud_sql_proxy | grep -v grep | awk '{print $2}')
}

case "$ACTION" in
  name)
    echo "cloud-sql-proxy";;
    myorder)
      echo 1;;
  status)
    if isRunning; then
      echo "running"
    else
      echo "stopped"
    fi;;
  start)
    startProxy;;
  stop)
    stopProxy;;
  restart)
    stopProxy
    sleep 1
    startProxy;;
  connect-check)
    exit 0;;
  connect)
    # To avoid double-checks, the script dose not check is-running.
    TZ=`date +%z`
    TZ=`echo ${TZ: 0: 3}:${TZ: -2}`
    # TODO: libray-ize and use 'isReceivingPipe' or even 'isInPipe' (suppress if piping in or out?)
    # test -t 0 && echo "Setting time zone: $TZ"
    # mysql -h127.0.0.1 "${CLOUDSQL_DB}" --init-command 'SET time_zone="'$TZ'"'
    psql "host=127.0.0.1 port=$PORT sslmode=disable dbname='$CLOUDSQL_DB' user='$CLOUDSQL_USER' password='$CLOUDSQL_PASSWORD'"
    ;;
  dump-check)
    exit 0;;
  dump)
    # currently, we only support a data-only dump
    # TODO: in future, make this a simple dump and take options as args past the action?
    mysqldump -h127.0.0.1 --skip-triggers --no-create-info --compatible=ansi --compact --complete-insert --single-transaction --ignore-table="${CLOUDSQL_DB}.catalystdb" "${CLOUDSQL_DB}";;
  param-default)
    ENV_PURPOSE="${1:-}"
    shift || (echo "Missing 'environment purpose' for 'param-default'." >&2; exit 1)
    PARAM_NAME="${1:-}"
    shift || (echo "Missing 'param name' for 'param-default'." >&2; exit 1)
    case "$ENV_PURPOSE" in
      dev|test)
        case "$PARAM_NAME" in
          CLOUDSQL_CONNECTION_NAME)
            echo '127.0.0.1:$PORT';;
          CLOUDSQL_CONNECTION_PROT)
            echo 'tcp';;
          *)
            echo '';;
        esac;;
      production|pre-production)
        case "$PARAM_NAME" in
          CLOUDSQL_CONNECTION_PROT)
            echo 'cloudsql';;
          *)
            echo '';;
        esac;;
      *)
        echo "Unknown environment purpose: '$ENV_PURPOSE'." >&2
        exit 1;;
    esac;;
  *)
    # TODO: library-ize and use 'echoerrandexit'
    echo "Unknown action '${ACTION}'." >&2
    exit 1;;
esac
