#!/bin/bash

SERVICE_NAME=robot-server.service
SCRIPTS_DIR=$PWD/v-arm/vr-controlled-robot-arm/raspberry/scripts
EXEC_DIR=$PWD/v-arm/server-build

cd $EXEC_DIR

systemctl stop $SERVICE_NAME
rm exec
mv temp/exec ./

systemctl start $SERVICE_NAME
journalctl -u $SERVICE_NAME -f
