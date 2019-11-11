#!/bin/bash

rm -fr            \
    $(pwd)/var    \
    $(pwd)/etc

docker container run                                    \
    --rm                                                \
    -v $(pwd)/var/lib/ceph:/var/lib/ceph                \
    -v $(pwd)/etc/ceph:/etc/ceph                        \
    -e CLUSTER_NAME=clusterfoo                          \
    -e MON_INIT_MEM_IDS=a                               \
    -e MON_INIT_MEM_ADDRS=172.21.0.10:6789              \
    -e PUBLIC_NETWORK=172.21.0.0/16                     \
    ceph:14.2.4                                         \
    ceph_cluster_init


#    -e CLUSTER_NETWORK=172.22.0.0/16                    \
