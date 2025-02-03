#!/bin/bash
JAR=note-prod-2024.04.jar
PROCESS=`jps | grep $JAR | awk '{print $1}'`
if [ -z "$PROCESS" ]; then
    java -Dfile.encoding=UTF-8 -Xss4096K -Xms1G -Xmx1G -jar note-prod-2024.04.jar >/dev/null 2>&1 &
	#java -DNOTE_DATASOURCE_URL=C:\Users\xiangqian\Desktop\tmp\note\data\database.db -Dfile.encoding=UTF-8 -Xss4096K -Xms1G -Xmx1G -jar note-prod-2024.04.jar --server.port=58080 >/dev/null 2>&1 &
else
    echo "Process already exists"
fi
