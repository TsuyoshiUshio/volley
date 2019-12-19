#!/bin/bash

# volley should be in the PATH 
LOGDIR=/home/azureuser/volley/log
mkdir -p $LOGDIR
nohup volley server &> ${LOGDIR}/volley.log &
echo "vollery server started. See the log at ${LOGDIR}/volley.log ."