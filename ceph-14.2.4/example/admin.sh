#!/bin/bash

#health status
docker exec -it mon.a ceph --cluster clusterfoo -s

#get monitor map from cluster and print it directly
docker exec -it mon.a ceph --cluster clusterfoo mon dump

#get monitor map from cluster, save it as a file, and then print it
monmapfile=monmap
docker exec -it mon.a ceph --cluster clusterfoo mon getmap -o /etc/ceph/$monmapfile
docker run --rm -it -v $(pwd):/app ceph-14.2.4:base /usr/bin/monmaptool --print /app/etc/ceph/$monmapfile
