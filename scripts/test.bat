@echo off

setlocal enabledelayedexpansion

set "DIRECTORY="

:parse_args
if "%~1"=="-d" (
  set "DIRECTORY=%~2"
  shift
  shift
  goto parse_args
)

if not defined DIRECTORY (
  echo You need to specify the source directory with the -d flag
  exit /b 1
)

if not exist "%DIRECTORY%" (
  echo Directory %DIRECTORY% does not exist
  exit /b 1
)


cd /d "%DIRECTORY%" || exit /b 1
go test -v -covermode count ./...
