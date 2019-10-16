#!/usr/bin/env bash

# clears any remaining files from a previous chain
rm -rf $HOME/.lyd
rm -rf $HOME/.lycli

../bin/lyd init longy --chain-id longychain

PWD_FILE=passwords.txt

if [ $# -eq 0 ]
then
    echo "No arguments supplied"
else
    PWD_FILE=$1
fi

# creates accounts keys with passwords, adds them to the keys list
../bin/lycli keys add alice < $PWD_FILE
../bin/lycli keys add bob < $PWD_FILE

# Add 2 accounts, with coins to the genesis file
../bin/lyd add-genesis-account $(../bin/lycli keys show alice -a) 1000longy,100000000stake
../bin/lyd add-genesis-account $(../bin/lycli keys show bob -a) 1000longy,100000000stake

# Set the default master key
../bin/lyd set-genesis-service

# Generate the genesis attendees from the eventbrite api
../bin/lyd add-genesis-attendees

# Generate the genesis prizes from the event
../bin/lyd add-genesis-prizes

# Sets the redeem account
../bin/lyd add-redeem-account $(../bin/lycli keys show alice -a)

# Configure your CLI to eliminate need for chain-id flag
../bin/lycli config chain-id longychain
../bin/lycli config output json
../bin/lycli config indent true
../bin/lycli config trust-node true

../bin/lyd gentx --name alice < $PWD_FILE
../bin/lyd gentx --name bob < $PWD_FILE

# input the gentx into the genesis file so chain is aware of validators
../bin/lyd collect-gentxs

# validate genesis
../bin/lyd validate-genesis

../bin/lyd start --rpc.laddr "tcp://0.0.0.0:26657"
