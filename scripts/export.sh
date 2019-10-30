#!/usr/bin/env bash
echo "Exporting current state of app"
./bin/lyd export --for-zero-height > lyd_export.json

echo "Reseting the chain state to 0"
./bin/lyd unsafe-reset-all

echo "Moving exported state to genesis"
cp lyd_export.json ~/.lyd/config/genesis.json

echo "Restarting the chain"
./bin/lyd start

