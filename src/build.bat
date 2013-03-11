@echo off
set APP=AutoMakeLess.exe
set PWD=%cd%\..
set GOPATH=%PWD%\src\add-on;%PWD%

if exist %APP% del %APP%

echo "Building %APP%"
go build .

if exist src.exe (
    rename src.exe %APP%
    echo "OK"
)
