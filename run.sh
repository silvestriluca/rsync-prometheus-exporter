#!/bin/bash
FILE_TO_PROCESS=$1
case "$FILE_TO_PROCESS" in
  "") 
    go run exporter.go
    ;;
  *)
    go run exporter.go "$FILE_TO_PROCESS"
    ;;
esac