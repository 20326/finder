#!/bin/bash

# ip2region
wget -N -O ip2region.db https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.db

filesize=`ls -l ip2region.db | awk '{ print $5 }'`
echo "ip2region.db updated! size:$filesize"
# phone
