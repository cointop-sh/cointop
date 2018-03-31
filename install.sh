#!/bin/bash

if ! [ $(id -u) = 0 ]; then
   echo "Must run as root"
   exit 1
fi

curl https://github.com/miguelmota/cointop/raw/master/bin/cointop -o cointop
mv cointop /usr/local/bin/
