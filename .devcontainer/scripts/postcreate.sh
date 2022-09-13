#!/bin/sh

set -e

# Install make
sudo apt-get update -yq
sudo apt-get install -yq make

# Install act from https://github.com/nektos/act
curl -s https://raw.githubusercontent.com/nektos/act/master/install.sh | sh
