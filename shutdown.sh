#!/bin/bash
JAR=note-prod-2024.04.jar
PROCESS=`jps | grep $JAR | awk '{print $1}'`
if [ -z "$PROCESS" ]; then
    echo "Process not found"
else
    for pid in $PROCESS; do
        echo "Kill process [ $pid ]"
        kill -9 $pid
    done
fi
