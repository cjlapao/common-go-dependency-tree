@echo off

setlocal enabledelayedexpansion

set "FILE="

:parse_args
if "%~1"=="" goto check_file

if "%~1"=="-f" (
  set "FILE=%~2"
  shift
  shift
  goto parse_args
)

echo Invalid option: %~1
exit /b 1

:check_file
if "%FILE%"=="" (
  echo You need to specify the version file with the -f flag
  exit /b
)

set /p VERSION=<"%FILE%"
echo %VERSION%

endlocal
