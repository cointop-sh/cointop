#!/bin/bash

echo "downloading..."
curl -s "https://github.com/miguelmota/cointop/raw/master/bin/cointop" -o cointop
echo "installing..."
mv cointop /usr/local/bin/
echo "cointop installed."
