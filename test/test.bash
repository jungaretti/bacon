#!/bin/bash

set -e

# Find directory of script file
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
LIB_DIR="$SCRIPT_DIR/lib"
LIB_FILE="$LIB_DIR/lib.bash"

make

./bin/bacon ping
./bin/bacon deploy config.example.yml
./bin/bacon deploy --delete --create config.example.yml

source $LIB_FILE
