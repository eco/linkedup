# node
make init &

# rest
sleep 8
bin/lycli rest-server --chain-id longychain --trust-node --laddr "tcp://0.0.0.0:1317" &

# key service
sleep 8
bin/ks --localstack
