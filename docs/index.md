# ConsoleDot service in Go

This projects aims to cover all basic concepts needed for development of service for console.redhat.com.

It aims to be a basic API serving service with a database.

## Repository structure

All packages live under `/internal` to denote we do not intend to share these with other apps.

All binaries have a directory under `/cmd` these are the app entrypoints.