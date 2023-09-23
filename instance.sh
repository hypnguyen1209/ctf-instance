#!/bin/bash

DIR=$1
ID=$2
PORT=$3
TIMEOUT=$4'm'
NAME=$5

NEW_DOCKERCOMPOSE="$2_docker-compose.yml"

PATH_INSTANCE=$DIR/$NEW_DOCKERCOMPOSE

cp $DIR/docker-compose.yml $PATH_INSTANCE

sed -i 's/{{PORT}}/'$PORT'/g' $PATH_INSTANCE

sed -i 's/{{ID}}/'$ID'/g' $PATH_INSTANCE

docker-compose -p $ID'_'$NAME -f $PATH_INSTANCE up -d

(sleep $TIMEOUT && docker-compose -p $ID'_'$NAME -f $PATH_INSTANCE down 2>&1 &)