#!/bin/bash

docker run -d \
    --name gospa-data \
    -p 15432:5432 \
    -e POSTGRES_PASSWORD=1234 \
    postgres
