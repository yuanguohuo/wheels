#!/bin/bash

docker network create --driver=bridge --subnet=172.20.0.0/16 ceph-public
