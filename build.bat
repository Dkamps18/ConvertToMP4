@echo off
setlocal enabledelayedexpansion

IF exist build/ (rmdir /s/q build)
mkdir build

for /f "delims=" %%x in (version) do set version=%%x

SET GOOS=linux
SET GOARCH=amd64
go build -o build/converttomp4
tar -C build -a -cf build/converttomp4-%version%-Linux-%GOARCH%.tar.gz converttomp4
set idx=0
for /f %%F in ('certutil -hashfile build/converttomp4-%version%-Linux-%GOARCH%.tar.gz SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo converttomp4-%version%-Linux-%GOARCH%.tar.gz  %out1% >> build/checksum.txt
SET GOARCH=arm64
go build -o build/converttomp4
tar -C build -a -cf build/converttomp4-%version%-Linux-%GOARCH%.tar.gz converttomp4
set idx=0
for /f %%F in ('certutil -hashfile build/converttomp4-%version%-Linux-%GOARCH%.tar.gz SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo converttomp4-%version%-Linux-%GOARCH%.tar.gz  %out1% >> build/checksum.txt

SET GOOS=windows
SET GOARCH=amd64
go build -o build/converttomp4.exe
tar -C build -a -cf build/converttomp4-%version%-Windows-%GOARCH%.zip converttomp4.exe
set idx=0
for /f %%F in ('certutil -hashfile build/converttomp4-%version%-Windows-%GOARCH%.zip SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo converttomp4-%version%-Windows-%GOARCH%.zip   %out1% >> build/checksum.txt
SET GOARCH=arm64
go build -o build/converttomp4.exe
tar -C build -a -cf build/converttomp4-%version%-Windows-%GOARCH%.zip converttomp4.exe
set idx=0
for /f %%F in ('certutil -hashfile build/converttomp4-%version%-Windows-%GOARCH%.zip SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo converttomp4-%version%-Windows-%GOARCH%.zip   %out1% >> build/checksum.txt

SET GOOS=darwin
SET GOARCH=amd64
go build -o build/converttomp4
tar -C build -a -cf build/converttomp4-%version%-macOS-%GOARCH%.tar.gz converttomp4
set idx=0
for /f %%F in ('certutil -hashfile build/converttomp4-%version%-macOS-%GOARCH%.tar.gz SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo converttomp4-%version%-macOS-%GOARCH%.tar.gz  %out1% >> build/checksum.txt
SET GOARCH=arm64
go build -o build/converttomp4
tar -C build -a -cf build/converttomp4-%version%-macOS-%GOARCH%.tar.gz converttomp4
set idx=0
for /f %%F in ('certutil -hashfile build/converttomp4-%version%-macOS-%GOARCH%.tar.gz SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo converttomp4-%version%-macOS-%GOARCH%.tar.gz  %out1% >> build/checksum.txt

rm build/converttomp4
rm build/converttomp4.exe