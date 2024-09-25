#!/bin/bash

set -e
set -o xtrace

# Function to check if the port is ready
wait_for_port() {
  local port=$1
  while ! nc -z localhost $port; do
    sleep 0.1
  done
}

repeat_cmd_until() {
  local cmd=$1
  local condition=$2
  local max_wait_secs=$3
  local interval_secs=2
  local start_time=$(date +%s)
  local output

  while true; do

    current_time=$(date +%s)
    if (( (current_time - start_time) > max_wait_secs )); then
      echo "Waited for expression "$1" to satisfy condition "$2" for $max_wait_secs seconds without luck. Returning with error."
      return 1
    fi

    output=$(eval $cmd)

    if [ $output $condition ]; then
      break
    else
      sleep $interval_secs
    fi
  done
}
