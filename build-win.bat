@echo off
chcp 65001 >nul
setlocal
cd /d "%~dp0"
go test ./... -count=1
if errorlevel 1 exit /b 1
call npm run build --prefix frontend
if errorlevel 1 exit /b 1
if not exist "产物" mkdir "产物"
set "GOOS=darwin"
set "GOARCH=amd64"
set "CGO_ENABLED=0"
go build -ldflags="-s -w" -o "产物\血棘苦修会秘典_macOS_amd64" .
if errorlevel 1 exit /b 1
set "GOOS=windows"
go build -ldflags="-s -w" -o "产物\血棘苦修会秘典.exe" .
if errorlevel 1 exit /b 1
endlocal
