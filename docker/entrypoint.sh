#!/bin/sh
# Starts the core. The data folder is the volume /var/lib/akad-vpn.
# The first start opens the panel and the proxy to the network
# (akad-vpn-core -setup-server) and prints the logins.
set -eu

data=/var/lib/akad-vpn
mkdir -p "$data"
if [ ! -f "$data/config.json" ]; then
  akad-vpn-core -data-dir "$data" -setup-server=-
fi
exec akad-vpn-core -plain -data-dir "$data"
