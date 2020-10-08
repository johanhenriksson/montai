#!/bin/bash
echo "creating montai mount point $1"

pid=$(pgrep docker-init)
if [ -z $pid ]
then
    # we are on a normal linux system
    # just set up a mount point
    nsenter -t 1 -m -- mkdir -p $1 
else
    # we're running in docker for mac
    # create a mount point AND a tempfs shared mount
    nsenter -t $pid -m -- mkdir -p $1 
    nsenter -t $pid -m -- mount -t tmpfs none --make-shared $1
fi
