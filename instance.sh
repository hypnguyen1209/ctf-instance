#!/bin/bash

DIR=$1
ID=$2
PORT=$3
NEWDOCKERCOMPOSE="$2_docker-compose.yml"

cp $DIR/docker-compose.yml $DIR/$NEWDOCKERCOMPOSE

sed -i 's/{{PORT}}/'$PORT'/g' $DIR/$NEWDOCKERCOMPOSE

sed -i 's/{{ID}}/'$ID'/g' $DIR/$NEWDOCKERCOMPOSE

docker-compose -p 