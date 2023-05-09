#!/bin/bash

# 当前目录
curDir=$(cd $(dirname $0); pwd)
echo "curDir: ${curDir}"
echo

# 输出目录
outputDir=${curDir}/build
echo "outputDir: ${outputDir}"
echo

# 删除输出目录
rm -rf "${outputDir}"
echo "rm ${outputDir}"
echo

# 创建输出目录
mkdir -p "${outputDir}"
echo "mkdir ${outputDir}"
echo

# cp
cp -r data "${outputDir}/"
echo "cp data"
cp -r i18n "${outputDir}/"
echo "cp i18n"
cp -r static "${outputDir}/"
echo "cp static"
cp -r template "${outputDir}/"
echo "cp template"
echo

# pkgName
os=`go env GOOS`
arch=`go env GOARCH`
pkgName=note_${os}_${arch}
echo "pkgName: ${pkgName}"
echo

# pkgPath & build
pkgPath=${outputDir}/${pkgName}
cd ./src && go build -ldflags="-s -w" -o "${pkgPath}"
# \$ sudo apt install upx
#cd ./src && go build -ldflags="-s -w" -o "${pkgPath}" && upx -9 --brute "${pkgPath}"
#cd ./src && go build -ldflags="-s -w" -o "${pkgPath}" && upx "${pkgPath}"
echo "pkgPath: ${pkgPath}"
echo

# startup.sh
startupPath=${outputDir}/startup.sh
cat>${startupPath}<<EOF
#!/bin/bash
# startup.sh
# \$ chmod +x ${pkgName} startup.sh
#nohup ./${pkgName} >/dev/null 2>&1 &
./${pkgName}
EOF
chmod +x "${startupPath}"
echo "startupPath: ${startupPath}"
echo