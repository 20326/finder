#!/bin/bash

# ip2region
wget -N -O ip2region.db https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.db

filesize=`ls -l ip2region.db | awk '{ print $5 }'`
oldsize=`ls -l etc/ip2region.db | awk '{ print $5 }'`
minsize=$((1024*1024*4))

if [ $filesize -gt $minsize ]
then
    if [ $filesize -eq $oldsize ]
    then
        echo "warring: ip2region.db, The two files are the same"
        rm ip2region.db
    else 
        mv etc/ip2region.db etc/ip2region."`date +%Y-%m-%d_%H:%M`".db
        mv ip2region.db etc/ip2region.db
        echo "ip2region.db updated! size:$filesize"
    fi
    
else 
    echo "error: ip2region.db  size:$filesize"
fi

# phone
