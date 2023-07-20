#!/bin/bash

# a simple bash script that reads a letter as its typed and echos it back to the screen
while read -n2 c
do
    echo -n "$c"
done
