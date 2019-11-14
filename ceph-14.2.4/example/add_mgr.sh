#!/bin/bash

source ./conf.sh

docker container run                                    \
    -d                                                  \
    -i                                                  \
    -t                                                  \
    --privileged=true                                   \
    --network=ceph-public                               \
    --name=mgr.a                                        \
    --ip=172.20.0.13                                    \
    -v $(pwd)/var/lib/ceph:/var/lib/ceph                \
    -v $(pwd)/var/run/ceph:/var/run/ceph                \
    -v $(pwd)/etc/ceph:/etc/ceph                        \
    -e CLUSTER_NAME=${CLUSTER_NAME}                     \
    -e MGR_ID=a                                         \
    ceph-14.2.4:v1                                      \
    ceph_add_mgr
