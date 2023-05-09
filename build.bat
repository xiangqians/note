::@echo off
::@echo on

:: 当前目录
set curDir=%~dp0
echo curDir: %curDir%
echo

:: 输出目录
set outputDir=%curDir%build
echo outputDir: %outputDir%
echo

:: 删除输出目录
:: rd /?
rd /s /q %outputDir%
echo rd %outputDir%
echo

:: copy
xcopy data "%outputDir%/data" /s /e /h /i /y
echo "cp data"
xcopy i18n "%outputDir%/i18n" /s /e /h /i /y
echo "cp i18n"
xcopy static "%outputDir%/static" /s /e /h /i /y
echo "cp static"
xcopy template "%outputDir%/template" /s /e /h /i /y
echo "cp template"
echo

:: pkgName
for /F %%i in ('go env GOOS') do (set os=%%i)
for /F %%i in ('go env GOARCH') do (set arch=%%i)
set pkgName=auto_deploy_%os%_%arch%.exe
echo pkgName: %pkgName%
echo

:: pkgPath & build
set pkgPath="%outputDir%/%pkgName%"
cd ./src && go build -ldflags="-s -w" -o %pkgPath%
::cd ./src && go build -ldflags="-s -w" -o %pkgPath% && upx -9 --brute %pkgPath%
echo pkgPath: %pkgPath%
echo

:: startup.bat
set startupPath=%outputDir%/startup.bat
echo :: startup.bat > %startupPath%
echo %pkgName% >> %startupPath%
echo pause >> %startupPath%
echo startupPath: %startupPath%
echo

pause