#!/bin/bash

echo "downloading..."

BIN="cointop32"
if [ $(uname -m) == 'x86_64' ]; then
	BIN="cointop64"
fi

wget "https://github.com/miguelmota/cointop/raw/master/bin/$BIN" -O cointop

echo "installing..."
chmod +x cointop
mv cointop /usr/local/bin/
echo "cointop installed."
