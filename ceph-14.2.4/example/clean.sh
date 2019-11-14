#!/bin/bash

docker container stop mgr.a
docker container stop mon.a
docker container stop mon.b
docker container stop mon.c

docker container rm mgr.a
docker container rm mon.a
docker container rm mon.b
docker container rm mon.c

rm -fr etc var
