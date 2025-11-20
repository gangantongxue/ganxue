#!/bin/sh
# wait-for-services.sh

set -e

SERVICES=""
while [ "$1" != "--" ]; do
  SERVICES="$SERVICES $1"
  shift
done
shift
CMD="$@"

for SERVICE in $SERVICES; do
  host=$(echo $SERVICE | cut -d: -f1)
  port=$(echo $SERVICE | cut -d: -f2)
  echo "Waiting for $host:$port ..."

  while ! nc -z "$host" "$port"; do
    echo "  $host:$port not ready, waiting..."
    sleep 2
  done

  echo "  $host:$port is ready!"
done

exec $CMD
