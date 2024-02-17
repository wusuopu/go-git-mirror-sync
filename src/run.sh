#!/bin/sh

if [[ -z $GO_ENV ]]; then
  GO_ENV=production
fi
start_server () {
  ./goose --env $GO_ENV up
  ./app
}


if [ "$1" = 'start_server' ]; then
  start_server
else
  echo $@
  exec $@
fi
