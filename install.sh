#!/bin/bash

echo "downloading..."
curl -s "https://github.com/miguelmota/cointop/raw/master/bin/cointop" -o cointop
echo "installing..."
chmod +x cointop
mv cointop /usr/local/bin/
echo "cointop installed."
