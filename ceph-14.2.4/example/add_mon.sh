#!/bin/bash

source ./conf.sh

docker container run                                    \
    -d                                                  \
    -i                                                  \
    -t                                                  \
    --privileged=true                                   \
    --network=ceph-public                               \
    --name=mon.b                                        \
    --ip=172.20.0.11                                    \
    -v $(pwd)/var/lib/ceph:/var/lib/ceph                \
    -v $(pwd)/var/run/ceph:/var/run/ceph                \
    -v $(pwd)/etc/ceph:/etc/ceph                        \
    -e CLUSTER_NAME=${CLUSTER_NAME}                     \
    -e MON_ID=b                                         \
    ceph-14.2.4:v1                                      \
    ceph_add_mon


docker container run                                    \
    -d                                                  \
    -i                                                  \
    -t                                                  \
    --privileged=true                                   \
    --network=ceph-public                               \
    --name=mon.c                                        \
    --ip=172.20.0.12                                    \
    -v $(pwd)/var/lib/ceph:/var/lib/ceph                \
    -v $(pwd)/var/run/ceph:/var/run/ceph                \
    -v $(pwd)/etc/ceph:/etc/ceph                        \
    -e CLUSTER_NAME=${CLUSTER_NAME}                     \
    -e MON_ID=c                                         \
    ceph-14.2.4:v1                                      \
    ceph_add_mon
