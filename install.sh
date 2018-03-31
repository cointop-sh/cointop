#!/bin/bash

echo "downloading..."
wget "https://github.com/miguelmota/cointop/raw/master/bin/cointop" -O cointop
echo "installing..."
chmod +x cointop
mv cointop /usr/local/bin/
echo "cointop installed."
