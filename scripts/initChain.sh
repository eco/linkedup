#!/usr/bin/env bash

# clears any remaining files from a previous chain
rm -rf $HOME/.lyd
rm -rf $HOME/.lycli

bt init --chain-id testchain

PWD_FILE=passwords.txt

if [ $# -eq 0 ]
then
    echo "No arguments supplied"
else
    PWD_FILE=$1
fi

# creates accounts keys with passwords, adds them to the keys list
btcli keys add alice < $PWD_FILE
btcli keys add bob < $PWD_FILE

# Add 2 accounts, with coins to the genesis file
bt add-genesis-account $(lycli keys show alice -a) 1000gamecoin
bt add-genesis-account $(lycli keys show bob -a) 1000gamecoin

bt start