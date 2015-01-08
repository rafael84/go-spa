#!/bin/bash
set -e

which migrate > /dev/null || go get -v github.com/mattes/migrate

dir=$(dirname $0)

ip=$(boot2docker ip 2> /dev/null || echo '127.0.0.1')  
port=15432
password=1234

migrate -path $dir/migrations -url "postgres://postgres:${password}@${ip}:${port}/postgres?sslmode=disable" $@
