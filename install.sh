#!/bin/bash

if ! [ $(id -u) = 0 ]; then
   echo "Must run as root"
   exit 1
fi

echo "downloading..."
curl https://github.com/miguelmota/cointop/raw/master/bin/cointop -o cointop
echo "installing..."
mv cointop /usr/local/bin/
echo "cointop installed."
