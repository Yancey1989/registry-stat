#!/bin/bash
registry-stat \
  --container-name=$CONTAINER_NAME \
  --container-path=$CONTAINER_PATH \
  --record-file=$RECORD_FILE \
  --dbconnect=$DBCONNECT
