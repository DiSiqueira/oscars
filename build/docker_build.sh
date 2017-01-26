#!/bin/sh

docker stop $(docker ps -qa -f "name=oscars")
docker rm $(docker ps -qa -f "name=oscars")
docker rmi -f $(docker images | grep oscars|awk {'print $3'})

docker-compose --project-name "build" -f docker-compose.yml create

docker tag build_oscars_golang oscars_golang
