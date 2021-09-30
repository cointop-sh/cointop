#!/bin/bash

VERSION=$(curl --silent "https://api.github.com/repos/cointop-sh/cointop/releases/latest" | grep -Po --color=never '"tag_name": ".\K.*?(?=")')

OSNAME="linux"
if [[ $(uname) == 'Darwin' ]]; then
  OSNAME="darwin"
fi

(
  cd /tmp
  wget https://github.com/cointop-sh/cointop/releases/download/v${VERSION}/cointop_${VERSION}_${OSNAME}_amd64.tar.gz
  tar -xvzf cointop_${VERSION}_${OSNAME}_amd64.tar.gz cointop

  sudo mv cointop /usr/local/bin/cointop
  cointop --version
)
