#!/bin/bash

while true
do
  t="$(( ( RANDOM % 20 ) + 10 ))"
  mosquitto_pub -h mosquitto -t temperature -m "${t}"
  echo "Sent temperature ${t} on $(date)"
  sleep 3
done
