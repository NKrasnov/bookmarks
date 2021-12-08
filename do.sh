#!/bin/bash

case $1 in
    run)
        go run ./server/application/bm-api
    ;;
    build)
        go build -o . ./server/application/bm-api
    ;;
esac