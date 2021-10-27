#!/usr/bin/env bash

docker-compose stop && docker-compose rm -f
docker images | grep 'app-cdn-serv' | awk '{print $3}' | xargs docker rmi -f
docker images | grep '<none>' | awk '{print $3}' | xargs docker rmi -f

docker-compose build && docker-compose up -d
