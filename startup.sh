#!/bin/bash
java -Dfile.encoding=UTF-8 -Xss4096K -Xms1G -Xmx1G -jar note-prod-2024.04.jar >/dev/null 2>&1 &
