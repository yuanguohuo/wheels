#!/bin/bash

source ./conf.sh

#by default, the initial monitor (mon.a) doesn't enable msgr2, enable it here;
docker exec -it mon.a ceph --cluster ${CLUSTER_NAME} mon enable-msgr2 
