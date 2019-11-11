#!/bin/bash

cluster=clusterfoo
mon_id=a



docker container run                                    \
    -d                                                  \
    -i                                                  \
    -t                                                  \
    --pid=host                                          \
    --privileged=true                                   \
    --network=ceph-public                               \
    --name=cmon.a                                       \
    --ip=172.21.0.10                                    \
    -v $(pwd)/etc/ceph:/etc/ceph                        \
    -v $(pwd)/var/lib/ceph:/var/lib/ceph                \
    -v $(pwd)/var/run/ceph:/var/run/ceph                \
    -e CLUSTER_NAME=${cluster}                          \
    -e MON_ID=${mon_id}                                 \
    ceph:14.2.4                                         \
    ceph_mon
