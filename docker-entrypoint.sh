#!/bin/bash

set -e

start() {
    echo "run service"
    exec /go/run ${@}
}

case $1 in
start)
    $1 ${@:2}
    ;;
*)
    exec $1 ${@:2}
    ;;
esac