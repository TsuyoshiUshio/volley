#!/bin/bash

# jmeter should be in the path.

JMETER_DIR=/home/azureuser/jmeter
mkdir -p $JMETR_DIR
cd $JMETER_DIR
# TODO IF the size becomes big, we might need to consider Log rotation
nohup jmeter-server &> jmeter-server-stdin-out.log 
echo "jmeter server started. See the log at $JMETER_DIR/jmeter-server-stdin-out.log and $JMETER_DIR/jmeter.log ."

