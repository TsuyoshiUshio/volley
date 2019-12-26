#!/bin/bash

# volley should be in the PATH 
LOGDIR=${HOME}/volley/log
mkdir -p $LOGDIR
nohup volley ss &> ${LOGDIR}/volley_slave.log &
echo "vollery slave server started. See the log at ${LOGDIR}/volley.log ."