@echo off

setlocal enabledelayedexpansion

set "PACKAGE_NAME="
set "DIRECTORY="
set "OUT_DIR=out"

:parse_args
if "%~1"=="-p" (
  set "PACKAGE_NAME=%~2"
  shift
  shift
  goto :parse_args
)
if "%~1"=="-d" (
  set "DIRECTORY=%~2"
  shift
  shift
  goto :parse_args
)

if not defined DIRECTORY (
  echo You need to specify the source directory with the -d flag
  exit /b 1
)

if not exist "%DIRECTORY%" (
  echo Directory %DIRECTORY% does not exist
  exit /b 1
)

if not exist "%OUT_DIR%" (
  mkdir "%OUT_DIR%"
)

if not exist "%OUT_DIR%\binaries" (
  mkdir "%OUT_DIR%\binaries"
)

cd "%DIRECTORY%" || exit /b 1
go build -o ..\%OUT_DIR%\binaries\%PACKAGE_NAME%
