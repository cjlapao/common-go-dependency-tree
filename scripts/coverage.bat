@echo off

setlocal

set "DIRECTORY="
set "COVERAGE_DIR=coverage"

:parse_args
if "%~1"=="" goto :check_directory
if "%~1"=="-d" (
  set "DIRECTORY=%~2"
  shift
  shift
  goto :parse_args
)
echo Invalid option: %~1
exit /b 1

:check_directory
if "%DIRECTORY%"=="" (
  echo You need to specify the source directory with the -d flag
  exit /b
)
if not exist "%DIRECTORY%" (
  echo Directory %DIRECTORY% does not exist
  exit /b
)

if not exist "%COVERAGE_DIR%" (
  mkdir "%COVERAGE_DIR%"
)

cd /d "%DIRECTORY%" || exit /b 1
go test -coverprofile coverage.txt -covermode count -v ./...
gocov convert coverage.txt | gocov-xml > ..\%COVERAGE_DIR%\cobertura-coverage.xml
del coverage.txt

endlocal
