#!/usr/bin/env bash
LYD=bin/lyd
LYCLI=bin/lycli

# clears any remaining files from a previous chain
rm -rf $HOME/.lyd
rm -rf $HOME/.lycli

$LYD init longy --chain-id longychain

PWD_FILE=passwords.txt

if [ $# -eq 0 ]
then
    echo "No arguments supplied"
else
    PWD_FILE=$1
fi

# creates accounts keys with passwords, adds them to the keys list
$LYCLI keys add alice < $PWD_FILE
$LYCLI keys add bob < $PWD_FILE
$LYCLI keys add redeemer < $PWD_FILE

# Add 2 accounts, with coins to the genesis file
$LYD add-genesis-account $($LYCLI keys show alice -a) 1000longy,100000000stake
$LYD add-genesis-account $($LYCLI keys show bob -a) 1000longy,100000000stake
$LYD add-genesis-account $($LYCLI keys show redeemer -a) 1000longy,100000000stake

# Set the default master and bonus key
$LYD set-genesis-key-service
$LYD set-genesis-bonus-service

# Generate the genesis attendees from the eventbrite api
$LYD add-genesis-attendees

# Generate the genesis prizes from the event
$LYD add-genesis-prizes

# Sets the consensus configurations file for the node to quicken block times
$LYD consensus-config

# Configure your CLI to eliminate need for chain-id flag
$LYCLI config chain-id longychain
$LYCLI config output json
$LYCLI config indent true
$LYCLI config trust-node true

$LYD gentx --name alice < $PWD_FILE
$LYD gentx --name bob < $PWD_FILE

# input the gentx into the genesis file so chain is aware of validators
$LYD collect-gentxs

# validate genesis
$LYD validate-genesis

$LYD start --rpc.laddr "tcp://0.0.0.0:26657"
