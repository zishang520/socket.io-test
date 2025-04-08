@echo OFF

set "args=%1"
pushd "%~dp0"
setlocal ENABLEDELAYEDEXPANSION
rem Set the GOPROXY environment variable
Set GOPROXY=https://goproxy.io,direct
Set DEBUG=*

if /i "%args%"=="default" goto %args%
if /i "%args%"=="deps" goto %args%
if /i "%args%"=="fmt" goto %args%
if /i "%args%"=="clean" goto %args%
if /i "%args%"=="test" goto %args%
if /i "%args%"=="run" goto %args%

goto default

:default
    GOTO :EOF

:deps
    CALL go mod tidy
    CALL go work sync
    CALL go work vendor
    GOTO :EOF

:fmt
    CALL go fmt -mod=mod ./...
    GOTO :EOF

:clean
    CALL go clean -mod=mod -v -r ./...
    GOTO :EOF

:test
    CALL go clean -testcache
    CALL go test -race -cover -covermode=atomic -mod=mod ./...
    GOTO :EOF

:run
    cls
    set "param=%2"
    if "%param%"=="" (
        set "param=./..."
    )
    echo run %param%
    call go run -race -gcflags "all=-N -l" -v %param%
    goto :EOF