#!/bin/bash

source ./conf.sh

rm -fr            \
    $(pwd)/var    \
    $(pwd)/etc

docker container run                                    \
    -d                                                  \
    -i                                                  \
    -t                                                  \
    --privileged=true                                   \
    --network=ceph-public                               \
    --name=mon.a                                        \
    --ip=172.20.0.10                                    \
    -v $(pwd)/var/lib/ceph:/var/lib/ceph                \
    -v $(pwd)/var/run/ceph:/var/run/ceph                \
    -v $(pwd)/etc/ceph:/etc/ceph                        \
    -e CLUSTER_NAME=${CLUSTER_NAME}                     \
    -e MON_INIT_MEM_ID=a                                \
    -e MON_INIT_MEM_ADDR=172.20.0.10:6789               \
    -e PUBLIC_NETWORK=172.20.0.0/16                     \
    ceph-14.2.4:v1                                      \
    ceph_cluster_create

#    -e CLUSTER_NETWORK=172.22.0.0/16                    \
