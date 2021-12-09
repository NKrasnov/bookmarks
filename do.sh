#!/bin/bash

case $1 in
    run)
        go run ./server/app/bm-api
    ;;
    build)
        go build -o . ./server/app/bm-api
    ;;
esac