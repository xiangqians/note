#!/bin/bash
#java -DNOTE_DATASOURCE=C:\Users\xiangqian\Desktop\tmp\note\data\database.db -Dfile.encoding=UTF-8 -Xss4096K -Xms1G -Xmx1G -jar note-prod-2024.04.jar --server.port=1520 >/dev/null 2>&1 &
java -Dfile.encoding=UTF-8 -Xss4096K -Xms1G -Xmx1G -jar note-prod-2024.04.jar >/dev/null 2>&1 &
