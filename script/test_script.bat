@echo off
setlocal enabledelayedexpansion

set output_file=output.txt
set /a counter=0

echo Writing to %output_file%...
echo.

:loop
set /a counter+=1
echo %time% - Write %counter% >> %output_file%
ping -n 2 127.0.0.1 > nul
if %counter%==10 goto :end
goto :loop

:end
echo Writing completed.