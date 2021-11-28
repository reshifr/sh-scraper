#!/bin/bash

find_command() { # package
  curl \
    --socks5 127.0.0.1:9050 \
    -s "https://packages.debian.org/bullseye/all/$@/filelist" \
    --compressed |
  recode html | grep -Eo '\/usr\/bin.*$' | sed -e 's/\/usr\/bin\///g'
}

export -f find_command
apt-cache pkgnames | parallel -j 192 find_command
