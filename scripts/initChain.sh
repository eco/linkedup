#!/usr/bin/env bash

# clears any remaining files from a previous chain
rm -rf $HOME/.lyd
rm -rf $HOME/.lycli

lyd init longy --chain-id longychain

PWD_FILE=passwords.txt

if [ $# -eq 0 ]
then
    echo "No arguments supplied"
else
    PWD_FILE=$1
fi

# creates accounts keys with passwords, adds them to the keys list
lycli keys add alice < $PWD_FILE
lycli keys add bob < $PWD_FILE

# Add 2 accounts, with coins to the genesis file
lyd add-genesis-account $(lycli keys show alice -a) 1000longy,100000000stake
lyd add-genesis-account $(lycli keys show bob -a) 1000longy,100000000stake

# Generate the genesis attendees from the eventbrite api
lyd add-genesis-attendees

# Set the default master key
lyd set-genesis-service

# Configure your CLI to eliminate need for chain-id flag
lycli config chain-id longychain
lycli config output json
lycli config indent true
lycli config trust-node true

lyd gentx --name alice < $PWD_FILE
lyd gentx --name bob < $PWD_FILE

# input the gentx into the genesis file so chain is aware of validators
lyd collect-gentxs

# validate genesis
lyd validate-genesis

lyd start
