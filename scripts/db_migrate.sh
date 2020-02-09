#!/bin/sh
cd docs/db
go-bindata -pkg migrations  .
mv bindata.go ../../pkg/model/migrations
gofmt -s -w ../../pkg/model/migrations/bindata.go
