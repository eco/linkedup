#!/usr/bin/env bash
LYD=../bin/lyd

echo "Exporting current state of app"
$LYD export --for-zero-height > lyd_export.json

echo "Reseting the chain state to 0"
$LYD unsafe-reset-all

echo "Moving exported state to genesis"
cp lyd_export.json ~/.lyd/config/genesis.json

echo "Restarting the chain"
$LYD start --rpc.laddr "tcp://0.0.0.0:26657"
