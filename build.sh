#!/bin/bash

# 当前目录
CURRENT_DIR=$(cd $(dirname $0); pwd)

# 输出目录
OUTPUT_DIR=${CURRENT_DIR}/build
echo "Output Dir: ${OUTPUT_DIR}"

# 清空输出目录
rm -rf "${OUTPUT_DIR}"

# 创建输出目录
mkdir -p "${OUTPUT_DIR}"

# 拷贝资源
# 拷贝note.ini文件
cp -r res/note.ini "${OUTPUT_DIR}/"
# 拷贝db文件
mkdir -p "${OUTPUT_DIR}/db"
cp -r res/db/database.db "${OUTPUT_DIR}/db/"

# 包名称
os=`go env GOOS`
arch=`go env GOARCH`
pkgName="note_${os}_${arch}"

# 包路径
pkgPath=${OUTPUT_DIR}/${pkgName}
echo $pkgPath

# 构建
cd ./src && go build -ldflags="-s -w" -o "${pkgPath}"
# \$ sudo apt install upx
#cd ./src && go build -ldflags="-s -w" -o "${pkgPath}" && upx -9 --brute "${pkgPath}"
#cd ./src && go build -ldflags="-s -w" -o "${pkgPath}" && upx "${pkgPath}"

# startup.sh
startupPath=${OUTPUT_DIR}/startup.sh
cat>${startupPath}<<EOF
#!/bin/bash
# startup.sh
#nohup ./${pkgName} >/dev/null 2>&1 &
./${pkgName}
EOF
chmod +x "${startupPath}"