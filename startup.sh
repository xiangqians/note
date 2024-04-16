#!/bin/bash
java -Dfile.encoding=UTF-8 -Xss4096K -Xms1G -Xmx1G -jar auto-deploy-prod-2024.02.jar >/dev/null 2>&1 &