#!/bin/bash
set -e

which migrate > /dev/null || go get -v github.com/mattes/migrate

# HOW TO CREATE USER AND DATABASE (assuming you're using psql):
#
#   $ psql template1
#   
#   template1=# create user gospa with password 'gospa';
#   template1=# create database gospa;
#   template1=# grant all privileges on database gospa to gospa;
#  
#   template1=# \c gospa;
#
#   gospa=# alter schema public owner to gospa;

dir=$(dirname $0)

. $dir/../server/.env

migrate -path $dir -url $DB_CONN_URL $@
