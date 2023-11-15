@echo off

REM 当前目录
set "CURRENT_DIR=%~dp0"

REM 输出目录
set "OUTPUT_DIR=%CURRENT_DIR%build"
echo Output Dir: %OUTPUT_DIR%

REM 清空输出目录
REM 判断输出目录是否已存在，如果已存在则删除
if exist %OUTPUT_DIR% (
    goto REM_DEL_OUTPUT_DIR
    REM 删除输出目录
    REM 删除一个目录：RMDIR [/S] [/Q] [drive:]path 或者 RD [/S] [/Q] [drive:]path
    REM /S 表示除目录本身外，还将删除指定目录下的所有子目录和文件，用于删除目录树。
    REM /Q 安静模式，带 /S 删除目录树时不要求确认。
    :REM_DEL_OUTPUT_DIR
    rd /s /q %OUTPUT_DIR%
)
REM 重新创建输出目录
md "%OUTPUT_DIR%"

REM 拷贝资源
REM 拷贝note.ini文件
copy "res\note.ini" "%OUTPUT_DIR%" > nul
REM 拷贝db文件
md "%OUTPUT_DIR%\db"
copy "res\db\database.db" "%OUTPUT_DIR%\db" > nul
REM 拷贝data目录
goto R
REM Xcopy命令格式：XCOPY source [destination]
REM 参数介绍：
REM source 指定要复制的文件
REM destination 指定新文件的位置和/或名称
REM /S 复制目录和子目录，除了空的
REM /E 复制目录和子目录，包括空的。与 /S /E 相同。可以用来修改 /T
REM /H 也复制隐藏和系统文件
REM /I 如果目标不存在，又在复制一个以上的文件，则假定目标一定是一个目录
REM /F 复制时显示完整的源和目标文件名
REM /Q 复制时不显示文件名
REM /T 创建目录结构，但不复制文件。不包括空目录或子目录。/T /E 包括空目录和子目录
REM /Y 禁止提示以确认改写一个现存目标文件
:R
xcopy "res\data" "%OUTPUT_DIR%\data" /s /e /h /i /y /f > nul

REM 包名称
for /F %%i in ('go env GOOS') do (set OS=%%i)
for /F %%i in ('go env GOARCH') do (set ARCH=%%i)
set "PKG_NAME=auto_deploy_%OS%_%ARCH%.exe"

REM 包路径
set "PKG_PATH=%OUTPUT_DIR%\%PKG_NAME%"

REM 构建
cd ./src && go build -ldflags="-s -w" -o %PKG_PATH%
::cd ./src && go build -ldflags="-s -w" -o %PKG_PATH% && upx -9 --brute %pkgPath%

REM startup.bat
set "STARTUP_PATH=%OUTPUT_DIR%\startup.bat"
echo %PKG_NAME% > %STARTUP_PATH%
echo pause >> %STARTUP_PATH%

pause