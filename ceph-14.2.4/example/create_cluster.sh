#!/bin/bash

rm -fr            \
    $(pwd)/var    \
    $(pwd)/etc

docker container run                                    \
    -d                                                  \
    -i                                                  \
    -t                                                  \
    --privileged=true                                   \
    --network=ceph-public                               \
    --ip=172.20.0.10                                    \
    -v $(pwd)/var/lib/ceph:/var/lib/ceph                \
    -v $(pwd)/etc/ceph:/etc/ceph                        \
    -e CLUSTER_NAME=clusterfoo                          \
    -e MON_INIT_MEM_ID=a                                \
    -e MON_INIT_MEM_ADDR=172.20.0.10:6789               \
    -e PUBLIC_NETWORK=172.20.0.0/16                     \
    ceph-14.2.4:v1                                      \
    ceph_cluster_create

#    -e CLUSTER_NETWORK=172.22.0.0/16                    \
