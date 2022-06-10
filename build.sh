#!/bin/bash
GOOS=linux GOARCH=amd64 go build
zip -qq -o postgresql-check-linux-amd64.zip postgresql-check
