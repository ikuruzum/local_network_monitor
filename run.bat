@echo off

cd /d "%~dp0"

if exist build rmdir /s /q build

if not exist build mkdir build

cd backend
go build -o ..\build\network_monitor_backend.exe
cd ..

cd frontend\local_network_monitor

flutter build windows --release

xcopy /E /I /Y build\windows\x64\runner\Release ..\..\build
