#!/bin/sh
set +x

pid=""

handle_sig_term() {
    echo "[entrypoint.sh] received SIGTERM, killing child."
    kill -TERM ${pid}
    wait ${pid}
}

handle_sig_int() {
    echo "[entrypoint.sh] received SIGINT, killing child."
    kill -INT ${pid}
    wait ${pid}
}

trap 'handle_sig_term' TERM
trap 'handle_sig_int' INT

./app -config=${CONFIG_FILE} & pid=$!
wait ${pid}
