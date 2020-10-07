#!/bin/bash
echo "creating shared mount point $1"
pid=$(pgrep docker-init)
echo "docker pid: $pid"
nsenter -t $pid -m -- mkdir -p $1 
nsenter -t $pid -m -- mount -t tmpfs none --make-shared $1
echo "ok"
