#!/bin/bash

VERSION=$(curl --silent "https://api.github.com/repos/miguelmota/cointop/releases/latest" | grep -Po --color=never '"tag_name": "\K.*?(?=")')

OSNAME="linux"
if [[ $(uname) == 'Darwin' ]]; then
  OSNAME="darwin"
fi

(
  cd /tmp
  wget https://github.com/miguelmota/cointop/releases/download/${VERSION}/cointop_${VERSION}_${OSNAME}_amd64.tar.gz
  tar -xvzf cointop_${VERSION}_${OSNAME}_amd64.tar.gz cointop

  sudo mv cointop /usr/local/bin/cointop
  cointop --version
)
