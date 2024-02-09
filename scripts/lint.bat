@echo off

set NAME=common-go-dependency-tree-linter
docker ps -a -f "name=%NAME%"

if %errorlevel% equ 1 (
  docker start "%NAME%" --attach
  exit /b %errorlevel%
)

docker run --name "%NAME%" -e RUN_LOCAL=true -e VALIDATE_ALL_CODEBASE=true -e VALIDATE_JSCPD=false -v .:/tmp/lint ghcr.io/super-linter/super-linter:slim-v5
set EXEC_CODE=%errorlevel%
docker rm "%NAME%"
exit /b %EXEC_CODE%
