#!/bin/bash

# To use this on Mac you need to brew install coreutils and then replace cut with gcut and numfmt with gnumfmt

git rev-list --objects --all \
| git cat-file --batch-check='%(objecttype) %(objectname) %(objectsize) %(rest)' \
| awk '/^blob/ {print substr($0,6)}' \
| sort --numeric-sort --key=2 \
| gcut --complement --characters=13-40 \
| gnumfmt --field=2 --to=iec-i --suffix=B --padding=7 --round=nearest