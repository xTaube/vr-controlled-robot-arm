#!/bin/bash

PROJECT_DIR=$PWD/v-arm/vr-controlled-robot-arm

cd $PROJECT_DIR

git checkout -f xtaube/remote-deploy-for-server > /dev/null
git pull > /dev/null
