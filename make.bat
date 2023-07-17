@echo OFF

set "args=%*"
pushd "%~dp0"
setlocal ENABLEDELAYEDEXPANSION
rem Set the GOPROXY environment variable
Set GOPROXY=https://proxy.golang.com.cn,direct
set http_proxy=socks5://127.0.0.1:1080
set https_proxy=%http_proxy%

if /i "%args%"=="update" goto %args%
if /i "%args%"=="install" goto %args%
if /i "%args%"=="all" goto %args%
if /i "%args%"=="engine.io" goto %args%
if /i "%args%"=="socket.io" goto %args%
if /i "%args%"=="init" goto %args%

goto DEFAULT_CASE
:update
    CALL go mod tidy
    CALL go mod vendor
    GOTO END_CASE
:install
    CALL go mod vendor -v
    GOTO END_CASE
:all
    echo ========================
    echo build
    set CGO_ENABLED=0
    CALL go build --mod=mod -race -ldflags "-s -w -extldflags \"-static\"" -o bin\main.exe main.go

    GOTO END_CASE
:engine.io
    set CGO_ENABLED=1
    CALL go build --mod=mod -race -o bin\engine.exe engine.io.go
    CALL bin\engine.exe
    GOTO END_CASE
:socket.io
    set CGO_ENABLED=1
    CALL go build --mod=mod -race -o bin\socket.exe socket.io.go
    CALL bin\socket.exe
    GOTO END_CASE
:init
    GOTO END_CASE
:DEFAULT_CASE
    set CGO_ENABLED=1
    CALL go build --mod=mod -race -o bin\main.exe main.go
    CALL bin\main.exe
    GOTO END_CASE
:END_CASE
    GOTO :EOF