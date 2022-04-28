@echo OFF

set "args=%*"
pushd "%~dp0"
setlocal ENABLEDELAYEDEXPANSION
set GOPATH="%~dp0vendor"
rem Set the GOPROXY environment variable
Set GOPROXY=https://goproxy.io,direct
set http_proxy=socks5://127.0.0.1:1080
set https_proxy=%http_proxy%

if /i "%args%"=="update" goto %args%
if /i "%args%"=="install" goto %args%
if /i "%args%"=="all" goto %args%
if /i "%args%"=="run" goto %args%
if /i "%args%"=="init" goto %args%

goto DEFAULT_CASE
:update
    if not exist vendor (
        CALL go mod vendor
    )
    CALL go mod tidy -v
    GOTO END_CASE
:install
    CALL go mod vendor -v
    GOTO END_CASE
:all
    echo ========================
    echo build
    set CGO_ENABLED=0
    CALL go build --mod=mod -ldflags "-s -w -extldflags \"-static\"" -o bin\main.exe main.go

    GOTO END_CASE
:run
    CALL go build --mod=mod -o bin\main.exe main.go
    CALL bin\main.exe
    GOTO END_CASE
:init
    GOTO END_CASE
:DEFAULT_CASE
    CALL go mod tidy
    set CGO_ENABLED=0
    CALL go build --mod=mod -ldflags "-s -w -extldflags \"-static\"" -o bin\main.exe main.go
    GOTO END_CASE
:END_CASE
    GOTO :EOF